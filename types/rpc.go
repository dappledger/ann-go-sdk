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

type ResultHealthInfo struct {
	Status int `json:"status"`
}

type ResultLastHeight struct {
	LastHeight int64 `json:"last_height"`
}

type ResultDialSeeds struct {
}

type ResultValidator struct {
	Address     []byte `json:"address"`
	PubKey      string `json:"pub_key"`
	VotingPower int64  `json:"voting_power"`
	Accum       int64  `json:"accum"`
	IsCA        bool   `json:"is_ca"`
}

type ResultValidators struct {
	BlockHeight int64        `json:"block_height"`
	Validators  []*ResultValidator `json:"validators"`
}

type ResultDumpConsensusState struct {
	RoundState      string   `json:"round_state"`
	PeerRoundStates []string `json:"peer_round_states"`
}

type ResultBroadcastTx struct {
	Code   CodeType `json:"code"`
	Data   []byte   `json:"data"`
	TxHash string   `json:"tx_hash"`
	Log    string   `json:"log"`
}

type ResultRequestSpecialOP struct {
	Code CodeType `json:"code"`
	Data []byte   `json:"data"`
	Log  string   `json:"log"`
}

type ResultBroadcastTxCommit struct {
	Code   CodeType `json:"code"`
	Data   []byte   `json:"data"`
	TxHash string   `json:"tx_hash"`
	Log    string   `json:"log"`
}

type ResultUnconfirmedTxs struct {
	N   int  `json:"n_txs"`
	Txs []Tx `json:"txs"`
}

type ResultNumArchivedBlocks struct {
	Num int64 `json:"num"`
}

type ResultNumLimitTx struct {
	Num uint64 `json:"num"`
}

type ResultInfo struct {
	Data             string `json:"data"`
	Version          string `json:"version"`
	LastBlockHeight  int64  `json:"last_block_height"`
	LastBlockAppHash []byte `json:"last_block_app_hash"`
}

type ResultQuery struct {
	Result Result `json:"result"`
}

type ResultRefuseList struct {
	Result []string `json:"result"`
}

type ResultBlock struct {
	BlockMeta *BlockMeta `json:"block_meta"`
	Block     *Block     `json:"block"`
}
