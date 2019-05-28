package sdk

import (
	"testing"
	"github.com/dappledger/AnnChain-go-sdk/abi"
	"strings"
	"github.com/dappledger/AnnChain-go-sdk/types"
	"encoding/json"
)

const (
	testEventABIStr = `[
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "_val",
				"type": "uint256"
			},
			{
				"indexed": false,
				"name": "val",
				"type": "uint256"
			}
		],
		"name": "SetVal",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "name",
				"type": "string"
			},
			{
				"indexed": false,
				"name": "sender",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "val",
				"type": "uint256"
			}
		],
		"name": "SetValByWho",
		"type": "event"
	}
]`;
	testEventReceiptStr = `{
	"PostState": "",
	"CumulativeGasUsed": 21464,
	"Bloom": "0x00000080000000000000000000000000000000000000000000000000000000000000010000001040000000000000000000004000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000820000000000000000100000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000080000000000000000000000000",
	"Logs": [
		{
			"address": "0x8490d1b29d30c33981436d477993039c18ba32d2",
			"topics": [
				"0xa44db7a20c59179d7fe5588218993c33539a9f6c9d94aee5e3df32b0b99105be",
				"0x00000000000000000000000000000000000000000000000000000000000000a8"
			],
			"data": "0x00000000000000000000000000000000000000000000000000000000000000a8",
			"blockNumber": "0x118c30",
			"transactionIndex": "0x0",
			"transactionHash": "0xf5d332f5794b6598a00b8e6eb50478a27de5896f63baaa03933fc95eb8f38f1b",
			"blockHash": "0xa083ad864eeb7dc4e78697ac7ab9ac6920c985e0bf55193c0914a55b0fe99c23",
			"logIndex": "0x0",
			"removed": false
		},
		{
			"address": "0x8490d1b29d30c33981436d477993039c18ba32d2",
			"topics": [
				"0x77cedbb85c40a0c2d19ca61f08c41f3d706d1233f5d2f29ac638685fdf1724b0"
			],
			"data": "0x000000000000000000000000000000000000000000000000000000000000006000000000000000000000000028112ca022224ae7757bcd559666be5340ff109a00000000000000000000000000000000000000000000000000000000000000a8000000000000000000000000000000000000000000000000000000000000002b61206e616d65207768696368206c656e67746820697320626967676572207468616e203332206279746573000000000000000000000000000000000000000000",
			"blockNumber": "0x118c30",
			"transactionIndex": "0x0",
			"transactionHash": "0xf5d332f5794b6598a00b8e6eb50478a27de5896f63baaa03933fc95eb8f38f1b",
			"blockHash": "0xa083ad864eeb7dc4e78697ac7ab9ac6920c985e0bf55193c0914a55b0fe99c23",
			"logIndex": "0x1",
			"removed": false
		}
	],
	"GasUsed": 21464,
	"Height": 5592
}`
)

func TestEventLog_String(t *testing.T) {
	eventsABI, _ := abi.JSON(strings.NewReader(testEventABIStr))
	eventsReceipt := new(types.ReceiptForStorage)
	if err := json.NewDecoder(strings.NewReader(testEventReceiptStr)).Decode(eventsReceipt); err != nil {
		t.Fatal(err.Error())
	}
	l := EventLog{(types.Receipt)(*eventsReceipt), eventsABI}

	expect := `SetVal(val:"168")`
	actual := l.Print("SetVal")
	if actual != expect {
		t.Fatal("event log unexpect", actual, expect)
	}

	expect = `SetValByWho(name:"a name which length is bigger than 32 bytes",sender:"0x28112ca022224ae7757bcd559666be5340ff109a",val:"168")`;
	actual = l.Print("SetValByWho")
	if actual != expect {
		t.Fatal("event log unexpect", actual, expect)
	}

	expect = `SetVal(val:"168")
SetValByWho(name:"a name which length is bigger than 32 bytes",sender:"0x28112ca022224ae7757bcd559666be5340ff109a",val:"168")`;
	actual = l.Print("SetVal", "SetValByWho")
	if actual != expect {
		t.Fatal("event log unexpect", actual, expect)
	}
}
