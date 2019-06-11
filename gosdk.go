package sdk

import (
	"strings"
)

type CyrptoType string

const (
	ZaCryptoType CyrptoType = "ZA"
)

type CommitType string

const (
	TypeSyn  CommitType = "syn"
	TypeAsyn CommitType = "asyn"
)

type GoSDK struct {
	rpcAddr    string
	cryptoType CyrptoType
	privkey    string
}

func (gs *GoSDK) Url() string {
	return gs.rpcAddr
}

func New(privKey, rpcAddr string, cryptoType CyrptoType) *GoSDK {

	if strings.Index(privKey, "0x") == 0 {
		privKey = privKey[2:]
	}

	return &GoSDK{
		rpcAddr,
		cryptoType,
		privKey,
	}
}

func (gs *GoSDK) JsonRPCCall(method string, params []byte, result interface{}) error {
	return gs.sendTxCall(method, params, result)
}

func (gs *GoSDK) Put(value []byte, typ CommitType) (string, error) {
	return gs.put(value, typ)
}

func (gs *GoSDK) Get(key string) ([]byte, error) {
	return gs.get(key)
}

func (gs *GoSDK) AccountCreate() (Account, error) {
	return gs.accountCreate()
}
