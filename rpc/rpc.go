// Copyright 2017 ZhongAn Information Technology Services Co.,Ltd.
//
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

package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	at "github.com/dappledger/AnnChain/angine/types"
	"github.com/dappledger/AnnChain/ann-module/lib/go-wire"
)

type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type RPCResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      string           `json:"id"`
	Result  *json.RawMessage `json:"result"`
	Error   *RPCError        `json:"error"`
}

type RPCError struct {
	Code    at.CodeType `json:"code"`
	Message string      `json:"message"`
}

func socketType(listenAddr string) string {
	socketType := "unix"
	if len(strings.Split(listenAddr, ":")) >= 2 {
		socketType = "tcp"
	}
	return socketType
}

func makeHTTPDialer(remoteAddr string) (string, func(string, string) (net.Conn, error)) {

	parts := strings.SplitN(remoteAddr, "://", 2)
	var protocol, address string
	if len(parts) != 2 {
		protocol = socketType(remoteAddr)
		address = remoteAddr
	} else {
		protocol, address = parts[0], parts[1]
	}

	trimmedAddress := strings.Replace(address, "/", ".", -1)
	return trimmedAddress, func(proto, addr string) (net.Conn, error) {
		return net.Dial(protocol, address)
	}
}

func makeHTTPClient(remoteAddr string) (string, *http.Client) {
	address, dialer := makeHTTPDialer(remoteAddr)
	return "http://" + address, &http.Client{
		Transport: &http.Transport{
			Dial: dialer,
		},
	}
}

type ClientJSONRPC struct {
	address string
	client  *http.Client
}

func NewClientJSONRPC(remote string) *ClientJSONRPC {
	address, client := makeHTTPClient(remote)
	return &ClientJSONRPC{
		address: address,
		client:  client,
	}
}

func (c *ClientJSONRPC) Call(method string, params []interface{}, result interface{}) ([]byte, at.CodeType, error) {
	request := RPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      "",
	}

	requestBytes := wire.JSONBytes(request)

	requestBuf := bytes.NewBuffer(requestBytes)
	httpResponse, err := c.client.Post(c.address, "text/json", requestBuf)
	if err != nil {
		return nil, at.CodeType_InternalError, err
	}
	defer httpResponse.Body.Close()
	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, at.CodeType_InternalError, err
	}
	return unmarshalResponseBytes(responseBytes, result)
}

func unmarshalResponseBytes(responseBytes []byte, result interface{}) ([]byte, at.CodeType, error) {

	var (
		err       error
		bytResult []byte
	)

	response := &RPCResponse{}

	if err = json.Unmarshal(responseBytes, response); err != nil {
		return nil, at.CodeType_DecodingError, errors.New(fmt.Sprintf("Error unmarshalling rpc response: %v", err))
	}

	if response.Error.Code != at.CodeType_OK {
		return nil, response.Error.Code, errors.New(response.Error.Message)
	}
	if response.Result == nil {
		return nil, at.CodeType_OK, nil
	}
	if bytResult, err = response.Result.MarshalJSON(); err != nil {
		return nil, at.CodeType_DecodingError, err
	}

	if result == nil {
		return bytResult, at.CodeType_OK, nil
	}

	if err = json.Unmarshal(bytResult, result); err != nil {
		return bytResult, at.CodeType_DecodingError, err
	}

	return bytResult, at.CodeType_OK, nil
}
