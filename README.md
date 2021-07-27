# Reeetip
Bot for tipping on Telegram

This bot was made for Viacoin only but I decided to support multiple assets.
What's special about this bot is that it doesn't need to be connected to a node.

Current supported coins are Bitcoin, Viacoin, Litecoin and Dash

### commands

```
tip - Send tip. Example /tip @RomanoRnr 2 via
deposit - Show your deposit address
balance - Show your balance
withdraw - Withdraw your funds
privkey - Show your private key
stats - Show tipping statistics
about - Get information about this bot
```

##### Minimum Recommended Specifications

- **Go 1.10 or 1.11**
* Linux


  Installation instructions can be found here: https://golang.org/doc/install.
  It is recommended to add `$GOPATH/bin` to your `PATH` at this point.

#### setup
``cd ~/go/src/gitlab.com/`` (create gitlab.com folder if doesn't exist)

``git clone https://github.com/viacoin/tipbot-telegram.git``

``reeetipbot``

``dep ensure``


dep is a dependency management tool for Go. It requires Go 1.9 or newer to compile.
https://github.com/golang/dep

#### config

Go into the config file and efit app.yml.example and rename it to app.yml

### Running the bot

You can run the bot by using the following command:

``go run main.go``

Or you can build the binary with the following command

``go build main.go``

This will produce a binary called "main" wich you can rename and for example upload
on the server with scp as example. With the compiled binary, your linux
server/vps does not need the dependency's. The server/VPS that will run the binary
does not need to have Golang installed at all.

The machine that will compile does need Golang installed and all dependency's (use dep).


----
