package sdk

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/dappledger/AnnChain-go-sdk/abi"
	"github.com/dappledger/AnnChain-go-sdk/common"
	"github.com/dappledger/AnnChain-go-sdk/crypto"
	"github.com/dappledger/AnnChain-go-sdk/ikhofi"
	"github.com/dappledger/AnnChain-go-sdk/types"
	"github.com/dappledger/AnnChain-go-sdk/utils"
)

const MAX_BATCH_PARAMS = 10000

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

func dealBatchParams(abiJson *abi.ABI, defaultMethod string, paramSlc []interface{}, dealFunc func(cdata []byte) error) error {
	if len(paramSlc) > MAX_BATCH_PARAMS {
		return fmt.Errorf("params should be less than %v", MAX_BATCH_PARAMS)
	}
	method := ""
	for i := range paramSlc {
		param, ok := paramSlc[i].([]interface{})
		if !ok {
			param, ok := paramSlc[i].(string)
			if !ok {
				return fmt.Errorf("not batch params")
			}
			method = param
			continue
		}
		if method == "" {
			method = defaultMethod
		}
		if method == "" {
			return fmt.Errorf("lack of method name")
		}
		calldata, err := toCallData(abiJson, method, param)
		if err != nil {
			return err
		}
		if err = dealFunc(calldata); err != nil {
			return err
		}
		method = ""
	}
	return nil
}

func PackCalldata(abiJson *abi.ABI, defaultMethod string, paramSlc []interface{}, batch bool) ([]byte, error) {
	var bb bytes.Buffer
	if !batch {
		return toCallData(abiJson, defaultMethod, paramSlc)
	}
	err := dealBatchParams(abiJson, defaultMethod, paramSlc, func(calldata []byte) error {
		utils.WriteBytes(&bb, calldata)
		return nil
	})
	alldata := bb.Bytes()
	alldata = append(types.BatchTag, alldata...)
	return alldata, err
}

func unpackResult(method string, abiDef abi.ABI, output string) (interface{}, error) {
	m, ok := abiDef.Methods[method]
	if !ok {
		return nil, fmt.Errorf("No such method")
	}
	if len(m.Outputs) == 1 {
		var result interface{}
		parsedData := ParseData(output)
		if err := abiDef.Unpack(&result, method, parsedData); err != nil {
			fmt.Println("error:", err)
			return nil, err
		}
		if strings.Index(m.Outputs[0].Type.String(), "bytes") == 0 {
			b, err := bytesN2Slice(result, m.Outputs[0].Type.Size)
			if err != nil {
				return nil, err
			}

			idx := 0
			for idx = 0; idx < len(b); idx++ {
				if b[idx] != 0 {
					break
				}
			}
			b = b[idx:]
			return fmt.Sprintf("0x%x", b), nil
		}
		return result, nil
	}

	d := ParseData(output)
	var result []interface{}
	if err := abiDef.Unpack(&result, method, d); err != nil {
		fmt.Println("fail to unpack outpus:", err)
		return nil, err
	}

	retVal := map[string]interface{}{}
	for i, output := range m.Outputs {
		var value interface{}
		if strings.Index(output.Type.String(), "bytes") == 0 {
			b, err := bytesN2Slice(result[i], m.Outputs[0].Type.Size)
			if err != nil {
				return nil, err
			}
			idx := 0
			for idx = 0; idx < len(b); idx++ {
				if b[idx] != 0 {
					break
				}
			}
			b = b[idx:]
			value = fmt.Sprintf("0x%x", b)
		} else {
			value = result[i]
		}
		if len(output.Name) == 0 {
			retVal[fmt.Sprintf("%v", i)] = value
		} else {
			retVal[output.Name] = value
		}
	}
	return retVal, nil
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
	if id == ikhofi.SystemContractId {
		if !(methodName == ikhofi.SystemDeployMethod ||
			methodName == ikhofi.SystemUpgradeMethod ||
			methodName == ikhofi.SystemQueryContractIdExits) {
			err = fmt.Errorf("Invalid system contract method: %s", methodName)
			return
		}
	}
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
	case strings.Index(typeName, "bool") == 0:
		if typeName == "bool" {
			return ParseBool(value)
		}
		return ParseBoolSlice(value, input.Type.Size)
	case strings.Index(typeName, "address") == 0:
		if typeName == "address" {
			return ParseAddress(value)
		}
		return ParseAddressSlice(value, input.Type.Size)
	case strings.Index(typeName, "uint8") == 0:
		if typeName == "uint8" {
			return ParseUint8(value)
		}
		return ParseUint8Slice(value, input.Type.Size)
	case strings.Index(typeName, "uint16") == 0:
		if typeName == "uint16" {
			return ParseUint16(value)
		}
		return ParseUint16Slice(value, input.Type.Size)
	case strings.Index(typeName, "uint32") == 0:
		if typeName == "uint32" {
			return ParseUint32(value)
		}
		return ParseUint32Slice(value, input.Type.Size)
	case strings.Index(typeName, "uint64") == 0:
		if typeName == "uint64" {
			return ParseUint64(value)
		}
		return ParseUint64Slice(value, input.Type.Size)
	case strings.Index(typeName, "int8") == 0:
		if typeName == "int8" {
			return ParseInt8(value)
		}
		return ParseInt8Slice(value, input.Type.Size)
	case strings.Index(typeName, "int16") == 0:
		if typeName == "int16" {
			return ParseInt16(value)
		}
		return ParseInt16Slice(value, input.Type.Size)
	case strings.Index(typeName, "int32") == 0:
		if typeName == "int32" {
			return ParseInt32(value)
		}
		return ParseInt32Slice(value, input.Type.Size)
	case strings.Index(typeName, "int64") == 0:
		if typeName == "int64" {
			return ParseInt64(value)
		}
		return ParseInt64Slice(value, input.Type.Size)
	case strings.Index(typeName, "uint256") == 0 ||
		strings.Index(typeName, "uint128") == 0 ||
		strings.Index(typeName, "int256") == 0 ||
		strings.Index(typeName, "int128") == 0:
		if typeName == "uint256" || typeName == "uint128" ||
			typeName == "int256" || typeName == "int128" {
			return ParseBigInt(value)
		}
		return ParseBigIntSlice(value, input.Type.Size)
	case strings.Index(typeName, "bytes8") == 0:
		if typeName == "bytes8" {
			return ParseBytesM(value, 8)
		}
		return ParseBytesMSlice(value, 8, input.Type.Size)
	case strings.Index(typeName, "bytes16") == 0:
		if typeName == "bytes16" {
			return ParseBytesM(value, 16)
		}
		return ParseBytesMSlice(value, 16, input.Type.Size)
	case strings.Index(typeName, "bytes32") == 0:
		if typeName == "bytes32" {
			return ParseBytesM(value, 32)
		}
		return ParseBytesMSlice(value, 32, input.Type.Size)
	case strings.Index(typeName, "bytes64") == 0:
		if typeName == "bytes64" {
			return ParseBytesM(value, 64)
		}
		return ParseBytesMSlice(value, 64, input.Type.Size)
	case strings.Index(typeName, "bytes") == 0:
		if typeName == "bytes" {
			return ParseBytes(value)
		}
	case typeName == "string":
		return ParseString(value)
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
