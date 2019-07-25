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

import (
	"time"
)

type BlockMeta struct {
	Hash        []byte        `json:"hash"`         // The block hash
	Header      *Header       `json:"header"`       // The block's Header
	PartsHeader PartSetHeader `json:"parts_header"` // The PartSetHeader, for transfer
}

type PartSetHeader struct {
	Total int    `json:"total"`
	Hash  []byte `json:"hash"`
}

type Block struct {
	*Header    `json:"header"`
	*Data      `json:"data"`
	LastCommit *Commit `json:"last_commit"`
}

//-----------------------------------------------------------------------------

type Header struct {
	ChainID         string    `json:"chain_id"`
	Height          int64     `json:"height"`
	Time            time.Time `json:"time"`
	NumTxs          int64     `json:"num_txs"` // XXX: Can we get rid of this?
	LastBlockID     BlockID   `json:"last_block_id"`
	LastCommitHash  []byte    `json:"last_commit_hash"` // commit from validators from the last block
	DataHash        []byte    `json:"data_hash"`        // transactions
	ValidatorsHash  []byte    `json:"validators_hash"`  // validators for the current block
	AppHash         []byte    `json:"app_hash"`         // state after txs from the previous block
	ReceiptsHash    []byte    `json:"recepits_hash"`    // recepits_hash from previous block
	ProposerAddress []byte    `json:"proposer_address"`
}

//-------------------------------------

// NOTE: Commit is empty for height 1, but never nil.
type Commit struct {
	// NOTE: The Precommits are in order of address to preserve the bonded ValidatorSet order.
	// Any peer with a block can gossip precommits by index with a peer without recalculating the
	// active ValidatorSet.
	BlockID    BlockID `json:"blockID"`
	Precommits []*Vote `json:"precommits"`
}

type Vote struct {
	ValidatorAddress []byte  `json:"validator_address"`
	ValidatorIndex   int     `json:"validator_index"`
	Height           int64   `json:"height"`
	Round            int64   `json:"round"`
	Type             byte    `json:"type"`
	BlockID          BlockID `json:"block_id"` // zero if vote is nil.
	Signature        string  `json:"signature"`
}

//-----------------------------------------------------------------------------

type Data struct {
	// Txs that will be applied by state @ block.Height+1.
	// NOTE: not all txs here are valid.  We're just agreeing on the order first.
	// This means that block.AppHash does not include these txs.
	Txs   Txs `json:"txs"`
	ExTxs Txs `json:"extxs"` // this is for all other txs which won't be delivered to app

	// Volatile
	hash []byte
}

//--------------------------------------------------------------------------------

type BlockID struct {
	Hash        []byte        `json:"hash"`
	PartsHeader PartSetHeader `json:"parts"`
}
