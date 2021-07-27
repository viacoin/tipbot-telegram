package user

import (
	"github.com/viacoin/viad/chaincfg/chainhash"
	"testing"
)

var x1, _ = chainhash.NewHashFromStr("2223")

//var utxo1 = UTXO {
//	x1,
//	1,
//	7,
//	true,
//	[]byte{},
//}
//
//var utxo2 = UTXO {
//	x1,
//	1,
//	4,
//	true,
//	[]byte{},
//}
//
//var utxo3 = UTXO {
//	x1,
//	1,
//	3,
//	true,
//	[]byte{},
//}
//
//var utxo4 = UTXO {
//	x1,
//	1,
//	8,
//	true,
//	[]byte{},
//}

var utxo1 = UTXO{
	x1,
	1,
	100000,
	true,
	[]byte{},
}

var utxo2 = UTXO{
	x1,
	1,
	100000,
	true,
	[]byte{},
}

var utxo3 = UTXO{
	x1,
	1,
	124860,
	true,
	[]byte{},
}

var utxo4 = UTXO{
	x1,
	1,
	3000000,
	true,
	[]byte{},
}

func TestMinimalRequiredUTXO(t *testing.T) {
	utxos := []*UTXO{&utxo1, &utxo2, &utxo3, &utxo4}
	getMinimalRequiredUTXO(int64(100000), utxos)
}
