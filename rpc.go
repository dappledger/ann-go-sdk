// Copyright © 2017 ZhongAn Technology
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sdk

import (
	"time"

	"github.com/dappledger/AnnChain-go-sdk/rpc"
	"github.com/dappledger/AnnChain-go-sdk/types"
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

func (gs *GoSDK) LastHeight() (int64, error) {
	rpcResult := &types.ResultLastHeight{}
	err := gs.sendTxCall("last_height", nil, rpcResult)
	if err != nil {
		return 0, err
	}
	return rpcResult.LastHeight, nil
}

func (gs *GoSDK) Block(height int64) (*types.ResultBlock, error) {
	clientJSON := rpc.NewClientJSONRPC(gs.rpcAddr)

	b := &types.ResultBlock{}

	if _, err := clientJSON.Call("block", []interface{}{height}, b); err != nil {
		return nil, err
	}

	return b, nil
}

func (c *GoSDK) IterBlock(from, to int64, iter func(prev, cur *types.ResultBlock) error) error {

	i := from

	var rprev *types.ResultBlock

	for i < to {
		last, err := c.LastHeight()
		if err != nil {
			return err
		}
		if last < i {
			time.Sleep(time.Second * 5)
			continue
		}

		for i <= last && i < to {
			b, err := c.Block(i)
			if err != nil {
				return err
			}
			if err := iter(rprev, b); err != nil {
				return err
			}
			rprev = b
			i++
		}
	}
	return nil
}
