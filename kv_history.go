package sdk

import (
	"encoding/binary"
	"fmt"

	"github.com/dappledger/ann-go-sdk/rlp"
	"github.com/dappledger/ann-go-sdk/types"
)

type ValueUpdateHistory struct {
	TxHash      []byte `json:"tx_hash"`
	BlockHeight uint64 `json:"block_height"`
	TimeStamp   uint64 `json:"time_stamp"`
	Value       []byte `json:"value"`
	TxIndex     uint32 `json:"tx_index"`
}

type ValueHistoryResult struct {
	Key                  []byte                `json:"key"`
	ValueUpdateHistories []*ValueUpdateHistory `json:"value_update_histories"`
	Total                uint32                `json:"total"`
}

func putUint32(i uint32) []byte {
	index := make([]byte, 4)
	binary.BigEndian.PutUint32(index, i)
	return index
}

func (gs *GoSDK) getKeyValueUpdateHistory(key []byte, pageNo uint32, pageSize uint32) (*ValueHistoryResult, error) {
	if pageNo == 0 {
		pageNo = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	query := append([]byte{types.QueryType_Key_Update_History}, putUint32(pageNo)...)
	query = append(query, putUint32(pageSize)...)
	query = append(query, key...)
	res := new(types.ResultQuery)
	err := gs.sendTxCall("query", query, res)
	if err != nil {
		return nil, err
	}
	if 0 != res.Result.Code {
		return nil, fmt.Errorf(string(res.Result.Log))
	}

	kvs := &ValueHistoryResult{
		ValueUpdateHistories: make([]*ValueUpdateHistory, 0),
	}

	if err := rlp.DecodeBytes(res.Result.Data, &kvs); err != nil {
		return nil, err
	}
	return kvs, nil
}
