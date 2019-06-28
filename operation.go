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

func (gs *GoSDK) accountCreate() (string, string) {
	return gs.mSigner.Gen()
}

func (gs *GoSDK) get(key string) ([]byte, error) {

	if strings.HasPrefix(key, "0x") {
		key = key[2:]
	}

	mParams := make(map[string]interface{}, 0)

	mParams["path"] = "GET"

	mParams["data"] = strings.ToLower(key)

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

	tx := trans.NewTransaction(address, big.NewInt(time.Now().UnixNano()), value)

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

func (gs *GoSDK) block(hashstr string) ([]string, int, error) {

	arryTxs := make([]string, 0)

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

	var txHashs []common.Hash

	err = rlp.DecodeBytes(result.Response.Value, &txHashs)
	if err != nil {
		return nil, 0, err
	}

	for _, txhash := range txHashs {
		arryTxs = append(arryTxs, txhash.Hex())
	}

	return arryTxs, len(arryTxs), nil
}
