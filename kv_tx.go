package sdk

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/dappledger/ann-go-sdk/common"
	"github.com/dappledger/ann-go-sdk/rlp"
	"github.com/dappledger/ann-go-sdk/types"
)

var KVTxType = []byte("kvTx-")

type KVTx struct {
	AccountBase
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

type KVResult struct {
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

func (gs *GoSDK) kvPut(kvTx *KVTx, funcType string) (string, error) {
	if kvTx.PrivKey == "" {
		return "", fmt.Errorf("account privkey is empty")
	}

	if strings.Index(kvTx.PrivKey, "0x") == 0 {
		kvTx.PrivKey = kvTx.PrivKey[2:]
	}

	privBytes := common.Hex2Bytes(kvTx.PrivKey)

	nonce := kvTx.Nonce
	if nonce == 0 {
		addrBytes, err := gs.getAddrBytes(privBytes)
		if err != nil {
			return "", err
		}
		nonce, err = gs.getNonce(common.Bytes2Hex(addrBytes))
		if err != nil {
			return "", err
		}
	}

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

	kvBytes, err := rlp.EncodeToBytes(&KVResult{Key: kvTx.Key, Value: kvTx.Value})
	if err != nil {
		return "", err
	}

	txdata := append(KVTxType, kvBytes...)

	tx := types.NewTransaction(nonce, common.Address{}, big.NewInt(0), gs.GasLimit(), big.NewInt(0), txdata)

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
	err = gs.sendTxCall(funcType, txBytes, rpcResult)
	if err != nil {
		return "", err
	}
	return rpcResult.TxHash, nil
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

func (gs *GoSDK) kvGetWithPrefix(prefix, lastKey []byte, limit uint32) ([]*KVResult, error) {
	params := struct {
		Prefix  []byte
		LastKey []byte
		Limit   uint32
	}{prefix, lastKey, limit}

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

func (gs *GoSDK) kvTxSigned(tx, funcType string) (string, error) {
	if strings.HasPrefix(tx, "0x") {
		tx = tx[2:]
	}

	txBytes := common.Hex2Bytes(tx)
	rpcResult := new(types.ResultBroadcastTxCommit)
	err := gs.sendTxCall(funcType, txBytes, rpcResult)
	if err != nil {
		return "", err
	}
	return rpcResult.TxHash, nil
}