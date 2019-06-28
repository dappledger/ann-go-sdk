package sdk

import (
	rpc "github.com/tendermint/tendermint/rpc/lib/client"
)

func (gs *GoSDK) sendTxCall(method string, params map[string]interface{}, result interface{}) error {

	clientJSON := rpc.NewJSONRPCClient(gs.rpcAddr)

	_, err := clientJSON.Call(method, params, result)

	return err
}
