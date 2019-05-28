package sdk

import (
	"strings"
)

type AccountBase struct {
	PrivKey string `json:"privkey"`
	Nonce   uint64 `json:"nonce"`
}

func (ab *AccountBase) ParseHex() {
	if strings.Index(ab.PrivKey, "0x") == 0 {
		ab.PrivKey = ab.PrivKey[2:]
	}
}

type ContractCreate struct {
	AccountBase
	Code   string        `json:"code"`
	ABI    string        `json:"abiDefinition"`
	Params []interface{} `json:"params"`
}

type ContractMethod struct {
	AccountBase
	Contract string        `json:"contract"`
	ABI      string        `json:"abiDefinition"`
	Method   string        `json:"method"`
	Params   []interface{} `json:"params"`
}

func (cm *ContractMethod) ParseHex() {
	if strings.Index(cm.Contract, "0x") == 0 {
		cm.Contract = cm.Contract[2:]
	}
	cm.AccountBase.ParseHex()
}

type QueryResult struct {
	Value string `json:"value"`
	Msg   string `json:"msg"`
}
