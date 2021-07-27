// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package bots

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/romanornr/CryptoTwitterTipBot/bcoins"
	"github.com/romanornr/CryptoTwitterTipBot/config"
	"github.com/romanornr/CryptoTwitterTipBot/database"
	"github.com/romanornr/CryptoTwitterTipBot/user"
	"github.com/spf13/viper"
	"github.com/viacoin/viad/chaincfg"
	"log"
	"strconv"
	"strings"
	"time"
)

func init() {
	config.GetViperConfig()
	database.Setup()
}

type TelegramBot struct {
	api string
	bot *tgbotapi.BotAPI
}

func NewTelegramBot(api string) *TelegramBot {
	bot, err := tgbotapi.NewBotAPI(api)
	if err != nil {
		log.Panicf("error new telegram bot: %s\n", err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)
	return &TelegramBot{
		api,
		bot,
	}
}

func (telegram TelegramBot) Update() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := telegram.bot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("error receiving messages: %s\n", err)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		command := telegram.commandParser(msg.Text, update)
		telegram.handler(command, &update)
	}
}

type Command struct {
	Command string
	Params  []string
	From    *tgbotapi.User
}

// receives message /p via 260
// command will be /p@botusername and params will be a slice of via and 260
func (telegram TelegramBot) commandParser(message string, update tgbotapi.Update) Command {
	var command string

	//add username to the command if message does not contain username
	//example /p via into /p@mewnbot
	botUsername := telegram.bot.Self.UserName
	command = strings.Split(message, " ")[0]
	if !strings.Contains(message, botUsername) {
		command = strings.Split(message, " ")[0] + "@" + botUsername
	}
	temp := strings.TrimLeft(message, botUsername)
	params := strings.Split(temp, " ")[1:]

	return Command{Command: command, Params: params, From: update.Message.From}
}

func (telegram TelegramBot) handler(command Command, update *tgbotapi.Update) {

	username := "@" + telegram.bot.Self.UserName
	m := map[string]interface{}{
		"/tip" + username:        telegram.Tip,
		"/deposit" + username:    telegram.Deposit,
		"/balance" + username:    telegram.getBalance,
		"/withdraw" + username:   telegram.Withdraw,
		"/privatekey" + username: telegram.receivePrivateKey,
		"/privkey" + username:    telegram.receivePrivateKey,
		"/about" + username:      telegram.about,
		"/stats" + username:      telegram.stats,
	}

	if _, ok := m[command.Command]; ok {
		m[command.Command].(func(command Command, update *tgbotapi.Update))(command, update)
	}
}

func (telegram TelegramBot) setNetwork(symbol string, update *tgbotapi.Update) (bcoins.Network, error) {
	coinsymbol := strings.ToLower(symbol)
	network, err := user.SelectNetwork(coinsymbol)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(viper.GetString("telegram_error_messages.altcoin_not_supported"), coinsymbol))
		telegram.bot.Send(msg)
		return bcoins.Network{}, err
	}
	chaincfg.Register(network.GetNetworkParams())
	return network, nil
}

// register the user who is using the bot
func RegisterTelegramSender(command Command) user.User {
	var user user.User
	user.Username = "@" + command.From.UserName
	user.Telegram_id = command.From.ID
	user.Social = "telegram"
	user.RegisteredDate = int64(time.Now().Unix())
	database.AddOrUpdateTelegramUser(user)
	user, _ = database.FindTelegramUserByUsername(user.Username)
	return user
}

// register user (not the person using the bot)
func RegisterTelegramReceiver(username string) user.User {
	if !strings.HasPrefix(username, "@") {
		return user.User{}
	}

	var receiver user.User
	receiver.Username = username
	receiver.Social = "telegram"
	receiver.RegisteredDate = int64(time.Now().Unix())
	database.AddOrUpdateTelegramUser(receiver)
	user, _ := database.FindTelegramUserByUsername(username)
	return user
}

// generate a unique public key for the user and send it in the chat
// deposit via
func (telegram TelegramBot) Deposit(command Command, update *tgbotapi.Update) {
	user := RegisterTelegramSender(command)

	if len(command.Params) < 1 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(viper.GetString("telegram_error_messages.deposit")))
		telegram.bot.Send(msg)
		return
	}

	network, err := telegram.setNetwork(command.Params[0], update)
	if err != nil {
		return
	}

	publicKey := user.PublicKey(network).String()
	message := fmt.Sprintf(viper.GetString("telegram_messages.deposit"), network.Name, network.Symbol, publicKey)
	warning := fmt.Sprintf(viper.GetString("telegram_messages.backup_privkey"), network.Symbol)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message+warning)
	coin, _ := bcoins.SelectCoin(command.Params[0])
	if coin.Binance == true {
		binanceShill := fmt.Sprintf("\n\n"+viper.GetString("shill_message.binance"), coin.Name)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, message+warning+binanceShill)
	}
	msg.DisableWebPagePreview = true

	telegram.bot.Send(msg)
}

// receive the private key in a private message.
// this private key can be imported/sweeped with third party software.
// command: /privekey via
func (telegram TelegramBot) receivePrivateKey(command Command, update *tgbotapi.Update) {
	if len(command.Params) < 1 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(viper.GetString("telegram_error_messages.privkey")))
		telegram.bot.Send(msg)
		return
	}

	user := RegisterTelegramSender(command)
	network, err := telegram.setNetwork(command.Params[0], update)
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(viper.GetString("telegram_messages.privkey_in_pm")))
	telegram.bot.Send(msg)
	msg = tgbotapi.NewMessage(int64(command.From.ID), fmt.Sprintf(viper.GetString("telegram_messages.privkey"), user.PrivateKeyWif(network).String()))
	telegram.bot.Send(msg)
}

// tip a user by the command /tip 1 via @Romanornr
func (telegram TelegramBot) Tip(command Command, update *tgbotapi.Update) {
	if len(command.Params) < 3 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(viper.GetString("telegram_error_messages.tipping")))
		telegram.bot.Send(msg)
		return
	}

	//check if the paramater contains a telegram username which starts with "@"
	if !strings.HasPrefix(command.Params[0], "@") {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(viper.GetString("telegram_error_messages.username"), command.Params[0]))
		telegram.bot.Send(msg)
		return
	}

	//check if the second param is float
	amount, err := strconv.ParseFloat(command.Params[1], 64)
	if err != nil {
		log.Printf("error float")
	}

	network, err := telegram.setNetwork(command.Params[2], update)
	if err != nil {
		return
	}

	receiver := RegisterTelegramReceiver(command.Params[0])
	sender := RegisterTelegramSender(command)

	txid, tx, err := sender.PayTo(network, receiver.PublicKey(network).String(), amount)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s!", err))
		telegram.bot.Send(msg)
		return
	}

	database.AddOrUpdateTotalTransferred(tx) // update total transferred amount in db

	insight, _ := user.GetInsightExplorer(network.Symbol)
	insightExplorer := fmt.Sprintf("%s/tx/%s", insight.Explorer, txid.Txid)
	message := fmt.Sprintf("Successfully tipped %s %s: %s", receiver.Username, strconv.FormatFloat(amount, 'f', -1, 64), insightExplorer)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	msg.ReplyToMessageID = update.Message.MessageID

	telegram.bot.Send(msg)
}

// withdraw funds by using the command /withdraw 2.5 VdMPvn7vUTSzbYjiMDs1jku9wAh1Ri2Y1A
// the bot will reply with a blockexplorer link to the transaction id.
func (telegram TelegramBot) Withdraw(command Command, update *tgbotapi.Update) {
	if len(command.Params) < 2 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(viper.GetString("telegram_error_messages.withdraw")))
		telegram.bot.Send(msg)
		return
	}

	network, err := telegram.setNetwork(command.Params[1], update)
	if err != nil {
		return
	}

	owner := RegisterTelegramSender(command)
	destinationAddress := command.Params[0]
	txid, _, err := owner.Withdraw(network, destinationAddress)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(viper.GetString("telegram_error_messages.something"), err))
		telegram.bot.Send(msg)
		return
	}
	insight, _ := user.GetInsightExplorer(network.Symbol)
	insightExplorer := fmt.Sprintf("%s/tx/%s", insight.Explorer, txid.Txid)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(viper.GetString("telegram_messages.withdrawal"), insightExplorer))

	telegram.bot.Send(msg)
}

// get balance stats by using the insight blockexplorer api
// example command: /balance via
// ouput: Balance: 27.33320144 VIA ($14.91)
func (telegram TelegramBot) getBalance(command Command, update *tgbotapi.Update) {

	if len(command.Params) < 1 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(viper.GetString("telegram_error_messages.balance")))
		telegram.bot.Send(msg)
		return
	}

	network, err := telegram.setNetwork(command.Params[0], update)
	if err != nil {
		return
	}

	user := RegisterTelegramSender(command)
	stats, err := user.AddressInfo(network)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s!", err))
		telegram.bot.Send(msg)
		return
	}

	totalUSDPrice := bcoins.GetTotalUSD(stats.Balance, network.Name)
	var message = fmt.Sprintf(viper.GetString("telegram_messages.balance"), strconv.FormatFloat(stats.Balance, 'f', -1, 64), strings.ToUpper(network.Symbol), totalUSDPrice)

	if stats.UnconfirmedBalance > 0 {
		message += fmt.Sprintf(viper.GetString("telegram_messages.unconfirmed_balance"), strconv.FormatFloat(stats.Balance, 'f', -1, 64), strconv.FormatFloat(stats.UnconfirmedBalance, 'f', -1, 64), strings.ToUpper(network.Symbol), totalUSDPrice)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	msg.ReplyToMessageID = update.Message.MessageID
	telegram.bot.Send(msg)
}

//get statistic in this format:
// 3 users have introduced themselves to me !
// Tipping statistics:
// dash: 0 dash ($0.00)
// litecoin: 0 ltc ($0.00)
// viacoin: 0.002 via ($0.00)
// bitcoin: 0 btc ($0.00)
func (telegram TelegramBot) stats(command Command, update *tgbotapi.Update) {
	coinStats := database.GetAllCoinStats()

	userCount := database.GetTotalUsers()
	var message = fmt.Sprintf(viper.GetString("telegram_messages.user_stats"), userCount)
	for i := 0; i < len(coinStats); i++ {
		message += fmt.Sprintf(viper.GetString("telegram_messages.coin_stats"), coinStats[i].Name, strconv.FormatFloat(coinStats[i].Amount, 'f', -1, 64), coinStats[i].Symbol, bcoins.GetTotalUSD(coinStats[i].Amount, coinStats[i].Name))
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	msg.ReplyToMessageID = update.Message.MessageID
	telegram.bot.Send(msg)
}

// show about info of the bot
func (telegram TelegramBot) about(command Command, update *tgbotapi.Update) {
	message := fmt.Sprintf("%s\n", viper.GetString("about.info") + viper.GetString("about.author") + viper.GetString("about.donate") + viper.GetString("about.btc_address") + viper.GetString("about.coins")+viper.GetString("about.warning"))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	msg.ReplyToMessageID = update.Message.MessageID
	telegram.bot.Send(msg)
}