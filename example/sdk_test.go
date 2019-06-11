package smoke

import (
	"testing"
	"time"

	"github.com/dappledger/AnnChain-go-sdk"
)

const (
	accPriv = "2c04b8ef19f5e69d05c4e715f5105d9177acb576b127aedb448957263e4d91e4"
	accAddr = "0xfedf9Aca577234cde08A2a9d934a6F9D9B0f557A"
)

var client *sdk.GoSDK

func init() {
	client = sdk.New(accPriv, "127.0.0.1:46657", sdk.ZaCryptoType)
}

func TestPutGet(t *testing.T) {

	hash, err := client.Put([]byte("blockdb-stressing-test"), sdk.TypeAsyn)

	t.Log(hash, err)

	time.Sleep(time.Second * 2)

	value, err := client.Get(hash)

	t.Log(string(value), err)
}
