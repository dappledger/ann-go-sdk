package sdk

import (
	"github.com/dappledger/AnnChain-go-sdk/rpc"
)

func (gs *GoSDK) sendTxCall(method string, params []byte, result interface{}) error {
	clientJSON := rpc.NewClientJSONRPC(gs.rpcAddr)
	var _params []interface{}
	if len(params) > 0 {
		_params = []interface{}{params}
	}
	_, err := clientJSON.Call(method, _params, result)
	return err
}
