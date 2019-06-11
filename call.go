package sdk

import (
	"strings"

	"github.com/dappledger/AnnChain-go-sdk/common"

	gtypes "github.com/dappledger/AnnChain/gemmill/types"
)

const (
	// angine takes query id from 0x01 to 0x2F
	QueryTxExecution = 0x01
)

func (gs *GoSDK) txSigned(tx string, isAsync bool) (string, error) {
	if strings.HasPrefix(tx, "0x") {
		tx = tx[2:]
	}
	txBys := common.Hex2Bytes(tx)
	rpcResult := new(gtypes.ResultBroadcastTxCommit)
	method := "broadcast_tx_commit"
	if isAsync {
		method = "broadcast_tx_async"
	}
	err := gs.sendTxCall(method, txBys, rpcResult)
	if err != nil {
		return "", err
	}
	hash := rpcResult.TxHash
	return hash, nil
}
