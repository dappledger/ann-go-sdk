package sdk

import (
	"crypto/ecdsa"
	"strings"

	"github.com/dappledger/AnnChain-go-sdk/ikhofi"
)

var (
	ikhofiSigner = ikhofi.DawnSigner{}
)

type ContractParam struct {
	ContractID string
	MethodName string
	Args       []string
	Privkey    *ecdsa.PrivateKey
	ByteCode   []byte
}

func getMethod(str string) (method string, args []string) {
	// get method from method strings
	index := strings.Index(str, "(")
	method = substr(str, 0, index)

	// get argument list from method strings
	argStr := substr(str, index+1, len(str)-index-2)
	args = []string{}
	if argStr != "" {
		args = strings.Split(argStr, ",")
	}
	for i := 0; i < len(args); i++ {
		arg := strings.Trim(args[i], "' ")
		args[i] = string(arg)
	}
	return method, args
}
