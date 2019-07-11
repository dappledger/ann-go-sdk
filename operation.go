package sdk

import (
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/dappledger/AnnChain/modules/common"
	"github.com/dappledger/AnnChain/modules/receipt"
	"github.com/dappledger/AnnChain/modules/rlp"
	trans "github.com/dappledger/AnnChain/modules/transaction"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type BlockTxHash struct {
	TxHash common.Hash
	Op     byte
}

func (gs *GoSDK) accountCreate() (string, string) {
	return gs.mSigner.Gen()
}

func (gs *GoSDK) get(privKey, key string) ([]byte, error) {

	if strings.HasPrefix(key, "0x") {
		key = key[2:]
	}

	if strings.HasPrefix(key, "0x") {
		privKey = privKey[2:]
	}

	address, err := gs.mSigner.PrivToAddress(privKey)
	if err != nil {
		return nil, err
	}

	tx := trans.NewTransaction(address, big.NewInt(time.Now().UnixNano()), common.Hex2Bytes(key), trans.OP_GET)

	sigTx, err := tx.SignTx(gs.mSigner, privKey)
	if err != nil {
		return nil, err
	}

	btx, err := rlp.EncodeToBytes(sigTx)
	if err != nil {
		return nil, err
	}

	mParams := make(map[string]interface{}, 0)

	mParams["path"] = "GET"

	mParams["data"] = common.Bytes2Hex(btx)

	var result ctypes.ResultABCIQuery

	err = gs.sendTxCall("abci_query", mParams, &result)
	if err != nil {
		return nil, err
	}

	if result.Response.Code != 0 {
		return nil, errors.New(result.Response.Log)
	}

	receipts := new(receipt.Receipt)

	if err = rlp.DecodeBytes(result.Response.Value, receipts); err != nil {
		return nil, err
	}

	return receipts.Value, nil

}

func (gs *GoSDK) put(privKey string, value []byte, typ CommitType) (string, error) {

	if strings.Index(privKey, "0x") == 0 {
		privKey = privKey[2:]
	}

	address, err := gs.mSigner.PrivToAddress(privKey)
	if err != nil {
		return "", err
	}

	tx := trans.NewTransaction(address, big.NewInt(time.Now().UnixNano()), value, trans.OP_PUT)

	sigTx, err := tx.SignTx(gs.mSigner, privKey)
	if err != nil {
		return "", err
	}

	txBytes, err := rlp.EncodeToBytes(sigTx)
	if err != nil {
		return "", err
	}

	mParams := make(map[string]interface{}, 0)

	mParams["tx"] = txBytes

	rpcCommitResult := new(ctypes.ResultBroadcastTxCommit)

	if typ == TypeSyn {
		err = gs.sendTxCall("broadcast_tx_commit", mParams, rpcCommitResult)
		if err != nil {
			return "", err
		}
		return sigTx.Hash(gs.mSigner).Hex(), nil
	}

	rpcAsynResult := new(ctypes.ResultBroadcastTx)

	err = gs.sendTxCall("broadcast_tx_async", mParams, rpcAsynResult)
	if err != nil {
		return "", err
	}

	return sigTx.Hash(gs.mSigner).Hex(), nil
}

func (gs *GoSDK) block(hashstr string) ([]BlockTxHash, int, error) {

	if strings.Index(hashstr, "0x") == 0 {
		hashstr = hashstr[2:]
	}

	mParams := make(map[string]interface{}, 0)

	mParams["path"] = "BLOCK"

	mParams["data"] = strings.ToLower(hashstr)

	var result ctypes.ResultABCIQuery

	err := gs.sendTxCall("abci_query", mParams, &result)
	if err != nil {
		return nil, 0, err
	}

	if result.Response.Code != 0 {
		return nil, 0, errors.New(result.Response.Log)
	}

	var txHashs []BlockTxHash

	err = rlp.DecodeBytes(result.Response.Value, &txHashs)
	if err != nil {
		return nil, 0, err
	}

	return txHashs, len(txHashs), nil
}

func (gs *GoSDK) get_log(hashstr string) (*receipt.Receipt, error) {

	if strings.HasPrefix(hashstr, "0x") {
		hashstr = hashstr[2:]
	}

	mParams := make(map[string]interface{}, 0)

	mParams["path"] = "LOG"

	mParams["data"] = strings.ToLower(hashstr)

	var result ctypes.ResultABCIQuery

	err := gs.sendTxCall("abci_query", mParams, &result)
	if err != nil {
		return nil, err
	}

	if result.Response.Code != 0 {
		return nil, errors.New(result.Response.Log)
	}

	receipts := new(receipt.Receipt)

	if err = rlp.DecodeBytes(result.Response.Value, receipts); err != nil {
		return nil, err
	}

	return receipts, nil
}
