package sdk

import (
	"math/big"

	"github.com/dappledger/AnnChain-go-sdk/types"
)

const GasLimit = 1000000000

type CyrptoType string

const (
	ZaCryptoType CyrptoType = "ZA"
)

type GoSDK struct {
	rpcAddr    string
	cryptoType CyrptoType
}

func (gs *GoSDK) Url() string {
	return gs.rpcAddr
}

func New(rpcAddr string, cryptoType CyrptoType) *GoSDK {
	return &GoSDK{
		rpcAddr,
		cryptoType,
	}
}

func (gs *GoSDK) JsonRPCCall(method string, params []byte, result interface{}) error {
	return gs.sendTxCall(method, params, result)
}

func (gs *GoSDK) AccountCreate() (Account, error) {
	return gs.accountCreate()
}

func (gs *GoSDK) Nonce(addr string) (uint64, error) {
	return gs.getNonce(addr)
}

func (gs *GoSDK) CheckHealth() (bool, error) {
	return gs.checkHealth()
}

func (gs *GoSDK) Block(hash string) ([]string, int, error) {
	return gs.block(hash)
}

func (gs *GoSDK) Receipt(hash string) (*types.ReceiptForStorage, error) {
	return gs.receipt(hash)
}

func (gs *GoSDK) Balance(addr string) (*big.Int, error) {
	return gs.balance(addr)
}

func (gs *GoSDK) ContractCreate(contract *ContractCreate) (map[string]interface{}, error) {
	return gs.contractCreate(contract)
}

func (gs *GoSDK) ContractCall(contractMethod *ContractMethod) (string, error) {
	return gs.contractCall(contractMethod, broadcast_tx_commit, false)
}

func (gs *GoSDK) ContractAsync(contractMethod *ContractMethod) (string, error) {
	return gs.contractCall(contractMethod, broadcast_tx_async, false)
}

func (gs *GoSDK) ContractRead(contractMethod *ContractMethod) (interface{}, error) {
	return gs.contractRead(contractMethod)
}

func (gs *GoSDK) Transaction(sendTx *Tx) (string, error) {
	return gs.sendTx(sendTx, "broadcast_tx_commit")
}

func (gs *GoSDK) TransactionAsync(sendTx *Tx) (string, error) {
	return gs.sendTx(sendTx, "broadcast_tx_async")
}

func (gs *GoSDK) TransactionPayLoad(txhash string) (string, error) {
	return gs.txPayLoad(txhash)
}

func (gs *GoSDK) TranscationSignature(tx string) (string, error) {
	return gs.txSigned(tx, false)
}

func (gs *GoSDK) TranscationSignatureAsync(tx string) (string, error) {
	return gs.txSigned(tx, true)
}
