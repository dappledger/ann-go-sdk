package sdk

import (
	"fmt"

	"github.com/dappledger/AnnChain-go-sdk/common"
	"github.com/dappledger/AnnChain-go-sdk/rlp"
	"github.com/dappledger/AnnChain-go-sdk/types"
)

type KVResult struct {
	K []byte
	V []byte
}

func (gs *GoSDK) kvGet(key []byte) ([]byte, error) {

	query := append([]byte{types.QueryType_Key}, key...)
	res := new(types.ResultQuery)
	err := gs.sendTxCall("query", query, res)
	if err != nil {
		return nil, err
	}
	if 0 != res.Result.Code {
		return nil, fmt.Errorf(string(res.Result.Log))
	}
	return res.Result.Data, nil
}

func (gs *GoSDK) kvGetWithPrefix(prefix, seeKey []byte, limit uint32) ([]*KVResult, error) {

	params := struct {
		Prefix []byte
		SeeKey []byte
		Limit  uint32
	}{prefix, seeKey, limit}

	bytParams, err := rlp.EncodeToBytes(params)
	if err != nil {
		return nil, err
	}
	query := append([]byte{types.QueryType_Key_Prefix}, bytParams...)
	res := new(types.ResultQuery)
	if err = gs.sendTxCall("query", query, res); err != nil {
		return nil, err
	}
	if 0 != res.Result.Code {
		return nil, fmt.Errorf(string(res.Result.Log))
	}
	kvs := make([]*KVResult, 0)

	if err := rlp.DecodeBytes(res.Result.Data, &kvs); err != nil {
		return nil, err
	}
	return kvs, nil
}

func (gs *GoSDK) kvPut(privKey string, nonce uint64, key, value []byte) (string, error) {
	privBytes := common.Hex2Bytes(privKey)
	addrBytes, err := gs.getAddrBytes(privBytes)
	if err != nil {
		return "", err
	}
	if nonce == 0 {
		nonce, err = gs.getNonce(common.Bytes2Hex(addrBytes))
		if err != nil {
			return "", err
		}
	}
	tx := types.NewKVTransaction(nonce, key, value)
	signer, sig, err := gs.signTx(privBytes, tx)
	if err != nil {
		return "", err
	}

	sigTx, err := tx.WithSignature(signer, sig)
	if err != nil {
		return "", err
	}
	txBytes, err := rlp.EncodeToBytes(sigTx)
	if err != nil {
		return "", err
	}

	rpcResult := new(types.ResultBroadcastTxCommit)
	err = gs.sendTxCall("broadcast_tx_commit", txBytes, rpcResult)
	if err != nil {
		return "", err
	}
	return rpcResult.TxHash, nil
}
