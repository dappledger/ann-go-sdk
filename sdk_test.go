package sdk

import (
	"testing"
)

const (
	accPriv = "48deaa73f328f38d5fcb29d076b2b639c8491f97d245fc22e95a86366687903a"
	accAddr = "28112ca022224ae7757bcd559666be5340ff109a"
)

var client *GoSDK

func init() {
	client = New(accPriv, "127.0.0.1:46657", ZaCryptoType)
}

func TestPutGet(t *testing.T) {

	hash, err := client.Put([]byte("myname"), TypeSyn)

	t.Log(hash, err)

	value, err := client.Get(hash)

	t.Log(value, err)
}
