package sdk

import (
	"github.com/dappledger/AnnChain/modules/receipt"
	"github.com/dappledger/AnnChain/modules/signer"
)

type CyrptoType string

const (
	Secp256K1 CyrptoType = "Secp256k1"
)

type CommitType string

const (
	TypeSyn  CommitType = "syn"
	TypeAsyn CommitType = "asyn"
)

type GoSDK struct {
	rpcAddr string
	mSigner signer.Signer
}

func (gs *GoSDK) Url() string {
	return gs.rpcAddr
}

func New(rpcAddr string, cryptoType CyrptoType) *GoSDK {
	switch cryptoType {
	case Secp256K1:
		return &GoSDK{rpcAddr: rpcAddr, mSigner: &signer.HomesteadSigner{}}
	}
	return nil
}

func (gs *GoSDK) Put(privKey string, value []byte, typ CommitType) (string, error) {
	return gs.put(privKey, value, typ)
}

func (gs *GoSDK) Get(privKey, key string) ([]byte, error) {
	return gs.get(privKey, key)
}

func (gs *GoSDK) AccountCreate() (string, string) {
	return gs.accountCreate()
}

func (gs *GoSDK) GetBlockTxs(blockHash string) ([]BlockTxHash, int, error) {
	return gs.block(blockHash)
}

func (gs *GoSDK) GetLog(key string) (*receipt.Receipt, error) {
	return gs.get_log(key)
}
