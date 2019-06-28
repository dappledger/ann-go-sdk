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
	client = New("127.0.0.1:26657", ZaCryptoType)
}

func TestPutGet(t *testing.T) {
	var (
		hash string
		err  error
	)

	t.Log(client.AccountCreate())

	hash, err = client.Put(accPriv, []byte("myname1"), TypeAsyn)
	hash, err = client.Put(accPriv, []byte("myname2"), TypeAsyn)
	hash, err = client.Put(accPriv, []byte("myname3"), TypeAsyn)
	hash, err = client.Put(accPriv, []byte("myname4"), TypeAsyn)

	t.Log(hash, err)

	time.Sleep(time.Second * 5)

	value, err := client.Get(hash)

	t.Log(string(value), err)

	txs, count, err := client.GetBlockTxs("e9bd0f4ece37595535b0bc721284107e8a822974c28e7a7ca247b38c1876864f")

	t.Log(count, err)

	for _, tx := range txs {
		value, err = client.Get(tx)
		t.Log(string(value), err)
	}
}
