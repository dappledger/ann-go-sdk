package sdk

import (
	"testing"

	"github.com/dappledger/AnnChain-go-sdk/common"
)

var (
	client *GoSDK
)

func init() {
	client = New("127.0.0.1:46657", ZaCryptoType)
}

func TestGetTransaction(t *testing.T) {
	result, tx, err := client.getTxByHash(common.Hex2Bytes("3cd32f6c96fc30f10e93136f533b01eab17eee5d742ea58c4df8007fc39a4bf4"))
	t.Log(result, tx, err)
}
