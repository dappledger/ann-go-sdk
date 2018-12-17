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

const (
	CREATE_ACCOUNT   = "create_account"
	PAYMENT          = "payment"
	MANAGE_DATA      = "manage_data"
	CREATE_CONTRACT  = "create_contract"
	EXECUTE_CONTRACT = "execute_contract"
	QUERY_CONTRACT   = "query_contract"
)

type CreateAccountParam struct {
	StartBalance string `json:"starting_balance"`
}

type PaymentParam struct {
	Amount string `json:"amount"`
}

type ManageDataValueParam struct {
	Value    string `json:"value"`
	Category string `json:"category"`
}

type ContractParam struct {
	PayLoad  string `json:"payload"`
	GasPrice string `json:"gas_price"`
	GasLimit string `json:"gas_limit"`
	Amount   string `json:"amount"`
}
