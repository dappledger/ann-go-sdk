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
	"crypto/ecdsa"
	"fmt"
	"strconv"
	"strings"

	"github.com/dappledger/ann-go-sdk/abi"
	"github.com/dappledger/ann-go-sdk/common"
	"github.com/dappledger/ann-go-sdk/crypto"
)

const MAX_BATCH_PARAMS = 10000

type ContractParam struct {
	ContractID string
	MethodName string
	Args       []string
	Privkey    *ecdsa.PrivateKey
	ByteCode   []byte
}

func substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

func getMethod(str string) (method string, args []string) {
	// get method from method strings
	index := strings.Index(str, "(")
	method = substr(str, 0, index)

	// get argument list from method strings
	argStr := substr(str, index+1, len(str)-index-2)
	args = []string{}
	if argStr != "" {
		args = strings.Split(argStr, ",")
	}
	for i := 0; i < len(args); i++ {
		arg := strings.Trim(args[i], "' ")
		args[i] = string(arg)
	}
	return method, args
}

func PackCalldata(abiJson *abi.ABI, defaultMethod string, paramSlc []interface{}) ([]byte, error) {
	return toCallData(abiJson, defaultMethod, paramSlc)
}

func unpackResult(method string, abiDef abi.ABI, output string) (interface{}, error) {
	m, ok := abiDef.Methods[method]
	if !ok {
		return nil, fmt.Errorf("no such method")
	}
	if len(m.Outputs) == 1 {
		var result interface{}
		parsedData := ParseData(output)
		if err := abiDef.Unpack(&result, method, parsedData); err != nil {
			fmt.Println("error:", err)
			return nil, err
		}
		return []interface{}{result}, nil
	}

	d := ParseData(output)
	result := make([]interface{}, m.Outputs.LengthNonIndexed())
	if err := abiDef.Unpack(&result, method, d); err != nil {
		fmt.Println("fail to unpack outputs:", err)
		return nil, err
	}
	return result, nil
}

func ParseData(data ...interface{}) (ret []byte) {
	for _, item := range data {
		switch t := item.(type) {
		case string:
			var str []byte
			if IsHex(t) {
				str = common.Hex2Bytes(t[2:])
			} else {
				str = []byte(t)
			}

			ret = append(ret, common.RightPadBytes(str, 32)...)
		case []byte:
			ret = append(ret, common.LeftPadBytes(t, 32)...)
		}
	}

	return
}

func IsHex(str string) bool {
	l := len(str)
	return l >= 4 && l%2 == 0 && str[0:2] == "0x"
}

func toCallData(abiJson *abi.ABI, method string, paramSlc []interface{}) ([]byte, error) {
	var calldata []byte
	args, err := parseArgs(method, abiJson, paramSlc)
	if err != nil {
		return nil, err
	}
	calldata, err = abiJson.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	return calldata, nil
}

func parseData(methodName string, abiDef *abi.ABI, params []interface{}) (string, error) {
	args, err := parseArgs(methodName, abiDef, params)
	if err != nil {
		return "", err
	}
	data, err := abiDef.Pack(methodName, args...)
	if err != nil {
		return "", err
	}

	var hexData string
	for _, b := range data {
		hexDataP := strconv.FormatInt(int64(b), 16)
		if len(hexDataP) == 1 {
			hexDataP = "0" + hexDataP
		}
		hexData += hexDataP
	}
	return hexData, nil
}

func parseArgs(methodName string, abiDef *abi.ABI, params []interface{}) ([]interface{}, error) {
	var method abi.Method
	if methodName == "" {
		method = abiDef.Constructor
	} else {
		var ok bool
		method, ok = abiDef.Methods[methodName]
		if !ok {
			return nil, fmt.Errorf("no such method")
		}
	}

	if params == nil {
		params = []interface{}{}
	}
	if len(params) != len(method.Inputs) {
		return nil, fmt.Errorf("unmatched params %x:%d", params, len(method.Inputs))
	}
	args := []interface{}{}

	for i := range params {
		a, err := ParseArg(method.Inputs[i], params[i])
		if err != nil {
			fmt.Println(fmt.Sprintf("fail to parse args %v into %s: %v ", params[i], method.Inputs[i].Name, err))
			return nil, err
		}
		args = append(args, a)
	}
	return args, nil
}

func SanitizeHex(hex string) string {
	return strings.TrimPrefix(strings.ToLower(hex), "0x")
}

func ParseParam(method, id, privkey string) (param ContractParam, err error) {
	if method == "" {
		err = fmt.Errorf("Required parameter method")
		return
	}
	if id == "" {
		err = fmt.Errorf("Required parameter id")
		return
	}
	if privkey == "" {
		err = fmt.Errorf("Required parameter privkey")
		return
	}
	privkey = SanitizeHex(privkey)
	ecdsaKey, errC := crypto.HexToECDSA(privkey)
	if errC != nil {
		err = errC
		return
	}

	methodName, args := getMethod(method)
	param = ContractParam{
		ContractID: id,
		MethodName: methodName,
		Args:       args,
		Privkey:    ecdsaKey,
	}
	return
}

func ParseArg(input abi.Argument, value interface{}) (interface{}, error) {
	typeName := input.Type.String()
	switch {
	case typeName == "bool":
		return ParseBool(value)
	case strings.HasPrefix(typeName, "bool"):
		return ParseBoolSlice(value)
	case typeName == "address":
		return ParseAddress(value)
	case strings.HasPrefix(typeName, "address"):
		return ParseAddressSlice(value)
	case typeName == "uint8":
		return ParseUint8(value)
	case strings.HasPrefix(typeName, "uint8"):
		return ParseUint8Slice(value)
	case typeName == "uint16":
		return ParseUint16(value)
	case strings.HasPrefix(typeName, "uint16"):
		return ParseUint16Slice(value)
	case typeName == "uint32":
		return ParseUint32(value)
	case strings.HasPrefix(typeName, "uint32"):
		return ParseUint32Slice(value)
	case typeName == "uint64":
		return ParseUint64(value)
	case strings.HasPrefix(typeName, "uint64"):
		return ParseUint64Slice(value)
	case typeName == "int8":
		return ParseInt8(value)
	case strings.HasPrefix(typeName, "int8"):
		return ParseInt8Slice(value)
	case typeName == "int16":
		return ParseInt16(value)
	case strings.HasPrefix(typeName, "int16"):
		return ParseInt16Slice(value)
	case typeName == "int32":
		return ParseInt32(value)
	case strings.HasPrefix(typeName, "int32"):
		return ParseInt32Slice(value)
	case typeName == "int64":
		return ParseInt64(value)
	case strings.HasPrefix(typeName, "int64"):
		return ParseInt64Slice(value)
	case strings.HasPrefix(typeName, "uint256") ||
		strings.HasPrefix(typeName, "uint128") ||
		strings.HasPrefix(typeName, "int256") ||
		strings.HasPrefix(typeName, "int128"):
		return ParseBigIntSlice(value)
	case typeName == "uint256" || typeName == "uint128" ||
		typeName == "int256" || typeName == "int128":
		return ParseBigInt(value)
	case typeName == "bytes8":
		return ParseBytesM(value, 8)
	case strings.HasPrefix(typeName, "bytes8"):
		return ParseBytesMSlice(value, 8)
	case typeName == "bytes16":
		return ParseBytesM(value, 16)
	case strings.HasPrefix(typeName, "bytes16"):
		return ParseBytesMSlice(value, 16)
	case typeName == "bytes32":
		return ParseBytesM(value, 32)
	case strings.HasPrefix(typeName, "bytes32"):
		return ParseBytesMSlice(value, 32)
	case typeName == "bytes64":
		return ParseBytesM(value, 64)
	case strings.HasPrefix(typeName, "bytes64"):
		return ParseBytesMSlice(value, 64)
	case typeName == "bytes":
		return ParseBytes(value)
	case typeName == "string":
		return ParseString(value)
	case strings.HasPrefix(typeName, "string"):
		return ParseStringSlice(value)
	}
	return nil, fmt.Errorf("type %v of %v is unsupported", typeName, input.Name)
}

func bytesN2Slice(value interface{}, m int) ([]byte, error) {
	switch m {
	case 0:
		v := value.([]byte)
		return v, nil
	case 8:
		v := value.([8]byte)
		return v[:], nil
	case 16:
		v := value.([16]byte)
		return v[:], nil
	case 32:
		v := value.([32]byte)
		return v[:], nil
	case 64:
		v := value.([64]byte)
		return v[:], nil
	}
	return nil, fmt.Errorf("type(bytes%d) not support", m)
}
