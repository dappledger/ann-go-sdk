package sdk

import (
	"strings"

	"github.com/dappledger/AnnChain-go-sdk/types"
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

type ResultTransaction struct {
	BlockHash        []byte `json:"block_hash"`
	BlockHeight      uint64 `json:"block_height"`
	TransactionIndex uint64 `json:"transaction_index"`
	RawTransaction   []byte `json:"raw_transaction"`
	Timestamp        uint64 `json:"timestamp"`
}

type ResultBlock struct {
	BlockMeta *types.BlockMeta `json:"block_meta"`
	Block     *types.Block     `json:"block"`
}
