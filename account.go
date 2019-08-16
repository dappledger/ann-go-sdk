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
	"github.com/dappledger/AnnChain-go-sdk/common"
	"github.com/dappledger/AnnChain-go-sdk/crypto"
)

type Account struct {
	Privkey string `json:"privkey"`
	Address string `json:"address"`
}

func (gs *GoSDK) accountCreate() (Account, error) {
	var account Account

	privkey, err := crypto.GenerateKey()
	if err != nil {
		return Account{}, err
	}

	address := crypto.PubkeyToAddress(privkey.PublicKey)

	account.Privkey = common.Bytes2Hex(crypto.FromECDSA(privkey))
	account.Address = address.Hex()

	return account, nil
}
