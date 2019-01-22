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

package annchain

import (
	"encoding/json"
	"strconv"

	"github.com/dappledger/AnnChain/genesis/types"
)

func NewCreateAccountTx(nonce uint64, basefee, memo, from, to, startBalance string) *types.Transaction {

	jspByte, _ := json.Marshal(CreateAccountParam{StartBalance: startBalance})

	return types.NewTransaction(strconv.FormatUint(nonce, 10), basefee, from, to, CREATE_ACCOUNT, memo, string(jspByte))
}

func NewRequestSpecialOPTx(isca bool, opcode uint8, validatorpub, from, rpcaddress, sigs string) *types.Transaction {

	jspByte, _ := json.Marshal(RequestSpecialOpParam{IsCA: isca, OpCode: opcode, ValidatorPub: validatorpub, RpcAddress: rpcaddress, Sigs: sigs})

	return types.NewTransaction("", "", "", "", REQUEST_SPECIAL_OP, "", string(jspByte))
}

func NewPaymentTx(nonce uint64, basefee, memo, from, to, amount string) *types.Transaction {

	jspByte, _ := json.Marshal(PaymentParam{Amount: amount})

	return types.NewTransaction(strconv.FormatUint(nonce, 10), basefee, from, to, PAYMENT, memo, string(jspByte))
}

func NewManageDataTx(nonce uint64, basefee, memo, from string, datas map[string]ManageDataValueParam) *types.Transaction {

	jspByte, _ := json.Marshal(datas)

	return types.NewTransaction(strconv.FormatUint(nonce, 10), basefee, from, "", MANAGE_DATA, memo, string(jspByte))
}

func NewCreateContractTx(nonce uint64, basefee, memo, from string, contract ContractParam) *types.Transaction {

	jspByte, _ := json.Marshal(contract)

	return types.NewTransaction(strconv.FormatUint(nonce, 10), basefee, from, "", CREATE_CONTRACT, memo, string(jspByte))
}

func NewExecuteContractTx(nonce uint64, basefee, memo, from, to string, contract ContractParam) *types.Transaction {

	jspByte, _ := json.Marshal(contract)

	return types.NewTransaction(strconv.FormatUint(nonce, 10), basefee, from, to, EXECUTE_CONTRACT, memo, string(jspByte))
}

func NewQueryContractTx(from, to string, contract ContractParam) *types.Transaction {

	jspByte, _ := json.Marshal(contract)

	return types.NewTransaction(strconv.FormatUint(0, 10), "0", from, to, QUERY_CONTRACT, "", string(jspByte))
}

func NewCreateContractParam(gasPrice, gasLimit, amount, byteCode, abis string, params []interface{}) (ContractParam, error) {

	var contractParam ContractParam

	payLoad, err := packCreateContractData(byteCode, abis, params)
	if err != nil {
		return contractParam, err
	}

	return ContractParam{
		GasLimit: gasLimit,
		GasPrice: gasPrice,
		Amount:   amount,
		PayLoad:  payLoad,
	}, nil
}

func NewExecuteContractParam(gasPrice, gasLimit, amount, funcName, abis string, params []interface{}) (ContractParam, error) {

	var contractParam ContractParam

	payLoad, err := packExecuteContractData(abis, funcName, params)
	if err != nil {
		return contractParam, err
	}

	return ContractParam{
		GasLimit: gasLimit,
		GasPrice: gasPrice,
		Amount:   amount,
		PayLoad:  payLoad,
	}, nil
}

func NewQueryContractParam(funcName, abis string, params []interface{}) (ContractParam, error) {

	var contractParam ContractParam

	payLoad, err := packExecuteContractData(abis, funcName, params)
	if err != nil {
		return contractParam, err
	}
	return ContractParam{
		PayLoad: payLoad,
	}, nil
}
