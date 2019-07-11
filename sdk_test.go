package sdk

import (
	"testing"
	"time"
)

const (
	accPriv = "48deaa73f328f38d5fcb29d076b2b639c8491f97d245fc22e95a86366687903a"
	accAddr = "28112ca022224ae7757bcd559666be5340ff109a"
)

var client *GoSDK

func init() {
	client = New("127.0.0.1:26657", Secp256K1)
}

func TestPutGet(t *testing.T) {
	var (
		hash string
		err  error
	)

	t.Log(client.AccountCreate())

	hash, err = client.Put(accPriv, []byte("myname1"), TypeAsyn)

	t.Log(hash, err)

	time.Sleep(time.Second * 2)

	value, err := client.Get(accPriv, hash)

	t.Log(string(value), err)

	txs, count, err := client.GetBlockTxs("9DE22D61EA4C167A1D20D0750781F4278B0E6B2661A16D700B2587D4E8C47463")

	t.Log(count, err)

	for _, tx := range txs {
		if tx.Op == 0x01 {
			value, err = client.Get(accPriv, tx.TxHash.Hex())
			t.Log(string(value), tx.Op, err)
		} else {
			rec, err := client.GetLog(tx.TxHash.Hex())
			t.Log(rec, tx.Op, err)
		}

	}
}
