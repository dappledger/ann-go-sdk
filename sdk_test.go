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
	client = New("127.0.0.1:46657", ZaCryptoType)
}

func TestPutGet(t *testing.T) {

	hash, err := client.Put(accPriv, []byte("myname"), TypeSyn)

	t.Log(hash, err)

	time.Sleep(time.Second * 5)

	value, err := client.Get(hash)

	t.Log(string(value), err)

	txs, count, err := client.GetBlockTxs("2D5622BA1C5E65AAD725D8BCF0DBD3D66F46D01E")

	t.Log(count, err)

	for _, tx := range txs {
		value, err = client.Get(tx)
		t.Log(string(value), err)
	}
}
