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

package types

// Volatile state for each Validator
// TODO: make non-volatile identity
// 	- Remove Accum - it can be computed, and now valset becomes identifying
type Validator struct {
	Address     []byte `json:"address"`
	PubKey      string `json:"pub_key"`
	VotingPower int64  `json:"voting_power"`
	Accum       int64  `json:"accum"`
	IsCA        bool   `json:"is_ca"`
}

type ValidatorAttr struct {
	PubKey []byte `json:"pubKey,omitempty"`
	Power  uint64 `json:"power,omitempty"`
	IsCA   bool   `json:"isCA,omitempty"`
}
