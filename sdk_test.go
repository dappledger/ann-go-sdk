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

	time.Sleep(time.Second)

	value, err := client.Get("0x75d101449a7b1983a389d47375f3d44448c32fd1947a5b0a78f67617a616eed8")

	t.Log(string(value), err)
}
