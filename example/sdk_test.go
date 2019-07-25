package smoke

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dappledger/AnnChain-go-sdk"
	"github.com/dappledger/AnnChain-go-sdk/abi"
	"github.com/dappledger/AnnChain-go-sdk/common"
	"github.com/dappledger/AnnChain-go-sdk/types"
	"github.com/stretchr/testify/assert"
)

var controlAbis string = `[
	{
		"constant": false,
		"inputs": [],
		"name": "destroyContract",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "_status",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [],
		"name": "changeReadOnly",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "contractName",
				"type": "string"
			},
			{
				"name": "newContract",
				"type": "address"
			}
		],
		"name": "changeContractAddress",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "_owner",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "queryStatus",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "contractName",
				"type": "string"
			}
		],
		"name": "queryContractAddress",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	}
]`

var tokenAbis string = `[
	{
		"constant": false,
		"inputs": [],
		"name": "destroyContract",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [],
		"name": "changeReadOnly",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "operationAccount",
				"type": "address"
			}
		],
		"name": "changeOperationAccount",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "tokenNo",
				"type": "string"
			},
			{
				"name": "dataType",
				"type": "uint8"
			},
			{
				"name": "newUintInfo",
				"type": "uint256"
			}
		],
		"name": "configureSingleUintTokenInfo",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "fromToken",
				"type": "string"
			},
			{
				"name": "toToken",
				"type": "string"
			}
		],
		"name": "comparTwoToken",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "tokenNo",
				"type": "string"
			}
		],
		"name": "queryTokenInfo",
		"outputs": [
			{
				"name": "",
				"type": "string"
			},
			{
				"name": "",
				"type": "string"
			},
			{
				"name": "",
				"type": "address"
			},
			{
				"name": "",
				"type": "address"
			},
			{
				"name": "",
				"type": "uint256"
			},
			{
				"name": "",
				"type": "uint256"
			},
			{
				"name": "",
				"type": "uint256"
			},
			{
				"name": "",
				"type": "uint256"
			},
			{
				"name": "",
				"type": "string"
			},
			{
				"name": "",
				"type": "uint256"
			},
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "operationAccount",
				"type": "address"
			}
		],
		"name": "setOperationAccount",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "account",
				"type": "address"
			},
			{
				"name": "isPermitted",
				"type": "bool"
			}
		],
		"name": "setAccessibleAccounts",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "from",
				"type": "string"
			},
			{
				"name": "parent",
				"type": "string"
			}
		],
		"name": "compareFromTokenTo",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "tokenNo",
				"type": "string"
			},
			{
				"name": "dataType",
				"type": "uint8"
			},
			{
				"name": "newAddressInfo",
				"type": "address"
			}
		],
		"name": "configureSingleAddressTokenInfo",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "tokenNo",
				"type": "string"
			},
			{
				"name": "dataType",
				"type": "uint8"
			}
		],
		"name": "querySingleAddressTokenInfo",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"name": "_accessibleAccounts",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "tokenNo",
				"type": "string"
			},
			{
				"name": "dateType",
				"type": "uint8"
			}
		],
		"name": "checkDateLine",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "symbol",
				"type": "string"
			},
			{
				"name": "token",
				"type": "string"
			}
		],
		"name": "chkSymbolToken",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "tokenNo",
				"type": "string"
			},
			{
				"name": "dataType",
				"type": "uint8"
			}
		],
		"name": "querySingleUintTokenInfo",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "queryCurrentAccounts",
		"outputs": [
			{
				"name": "",
				"type": "address[]"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "tokenNo",
				"type": "string"
			}
		],
		"name": "queryTokenIsExist",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "tokenNo",
				"type": "string"
			}
		],
		"name": "queryTokenFinOwner",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "fTokenNo",
				"type": "string"
			},
			{
				"name": "tTokenNo",
				"type": "string"
			}
		],
		"name": "checkFinOwner",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "tokenNo",
				"type": "string"
			},
			{
				"name": "dataType",
				"type": "uint256"
			}
		],
		"name": "querySingleBytesTokenInfo",
		"outputs": [
			{
				"name": "result",
				"type": "bytes"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "_owner",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "tokenNo",
				"type": "string"
			},
			{
				"name": "dataType",
				"type": "uint8"
			},
			{
				"name": "newStringInfo",
				"type": "string"
			}
		],
		"name": "configureSingleStringTokenInfo",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "symbol",
				"type": "string"
			},
			{
				"name": "tokenNo",
				"type": "string"
			},
			{
				"name": "accreditor",
				"type": "address"
			},
			{
				"name": "signer",
				"type": "address"
			},
			{
				"name": "creditBeginTime",
				"type": "uint256"
			},
			{
				"name": "creditDeadline",
				"type": "uint256"
			},
			{
				"name": "conversionRate",
				"type": "uint256"
			}
		],
		"name": "configureTokenInfo1",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "tokenNo",
				"type": "string"
			},
			{
				"name": "dataType",
				"type": "uint256"
			}
		],
		"name": "querySingleStringTokenInfo",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "from",
				"type": "string"
			},
			{
				"name": "to",
				"type": "string"
			}
		],
		"name": "chkSymbol",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "sysAccount",
				"type": "address"
			}
		],
		"name": "setSysAccount",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "tokenNo",
				"type": "string"
			},
			{
				"name": "signBeginTime",
				"type": "uint256"
			},
			{
				"name": "signDeadline",
				"type": "uint256"
			},
			{
				"name": "superiorTokenNo",
				"type": "string"
			},
			{
				"name": "generatedTime",
				"type": "uint256"
			}
		],
		"name": "configureTokenInfo2",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "_operationAccount",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	}
]`

var accountAbis string = `[
	{
		"constant": true,
		"inputs": [
			{
				"name": "account",
				"type": "address"
			},
			{
				"name": "symbol",
				"type": "string"
			}
		],
		"name": "queryPrivacy",
		"outputs": [
			{
				"name": "",
				"type": "address[]"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "symbol",
				"type": "string"
			},
			{
				"name": "accounts",
				"type": "address[]"
			}
		],
		"name": "setGlobalPrivacy",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [],
		"name": "destroyContract",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "_status",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [],
		"name": "changeReadOnly",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "operationAccount",
				"type": "address"
			}
		],
		"name": "changeOperationAccount",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "operationAccount",
				"type": "address"
			}
		],
		"name": "setOperationAccount",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "symbol",
				"type": "string"
			}
		],
		"name": "queryGlobalPrivacy",
		"outputs": [
			{
				"name": "",
				"type": "address[]"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "_unionAccount",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "account",
				"type": "address"
			},
			{
				"name": "isPermitted",
				"type": "bool"
			}
		],
		"name": "configureAccessibleAccounts",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "controlContract",
				"type": "address"
			}
		],
		"name": "initContracts",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"name": "_accessibleAccounts",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "account",
				"type": "address"
			},
			{
				"name": "tokenNo",
				"type": "string"
			},
			{
				"name": "unactivatedToken",
				"type": "uint256"
			},
			{
				"name": "standardToken",
				"type": "uint256"
			},
			{
				"name": "checkPendingToken",
				"type": "uint256"
			},
			{
				"name": "subscribedToken",
				"type": "uint256"
			},
			{
				"name": "frozenToken",
				"type": "uint256"
			}
		],
		"name": "configureAccountToken",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "from",
				"type": "uint256"
			},
			{
				"name": "count",
				"type": "uint256"
			}
		],
		"name": "queryAccouts",
		"outputs": [
			{
				"name": "",
				"type": "address[]"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "queryCurrentAccounts",
		"outputs": [
			{
				"name": "",
				"type": "address[]"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "account",
				"type": "address"
			},
			{
				"name": "symbol",
				"type": "string"
			},
			{
				"name": "tokenNo",
				"type": "string"
			}
		],
		"name": "queryAccountTokenInfo",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			},
			{
				"name": "",
				"type": "uint256"
			},
			{
				"name": "",
				"type": "uint256"
			},
			{
				"name": "",
				"type": "uint256"
			},
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "_owner",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "queryStatus",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "account",
				"type": "address"
			},
			{
				"name": "tokenNo",
				"type": "string"
			},
			{
				"name": "dataType",
				"type": "uint8"
			},
			{
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "configureSingleUintToken",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "account",
				"type": "address"
			},
			{
				"name": "fromIndex",
				"type": "uint256"
			}
		],
		"name": "queryAccountTokens",
		"outputs": [
			{
				"name": "",
				"type": "string"
			},
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"name": "_accounts",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "account",
				"type": "address"
			},
			{
				"name": "tokenNo",
				"type": "string"
			}
		],
		"name": "checkTokenExists",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "account",
				"type": "address"
			},
			{
				"name": "symbol",
				"type": "string"
			},
			{
				"name": "accounts",
				"type": "address[]"
			}
		],
		"name": "setPrivacy",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "sysAccount",
				"type": "address"
			}
		],
		"name": "setSysAccount",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "account",
				"type": "address"
			},
			{
				"name": "tokenNo",
				"type": "string"
			},
			{
				"name": "dataType",
				"type": "uint8"
			}
		],
		"name": "queryAccountTokenTypeAmount",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	}
]`

var transAbis string = `[
	{
		"constant": false,
		"inputs": [
			{
				"name": "voteId",
				"type": "uint256"
			},
			{
				"name": "account",
				"type": "address"
			},
			{
				"name": "symbol",
				"type": "string"
			},
			{
				"name": "tokenNo",
				"type": "string"
			},
			{
				"name": "rate",
				"type": "uint256"
			},
			{
				"name": "unactivate",
				"type": "uint256"
			},
			{
				"name": "standard",
				"type": "uint256"
			},
			{
				"name": "opType",
				"type": "uint256"
			}
		],
		"name": "publish",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [],
		"name": "destroyContract",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [],
		"name": "changeReadOnly",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "operationAccount",
				"type": "address"
			}
		],
		"name": "changeOperationAccount",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "operationAccount",
				"type": "address"
			}
		],
		"name": "setOperationAccount",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "controlContract",
				"type": "address"
			}
		],
		"name": "initContracts",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "model",
				"type": "uint256"
			},
			{
				"name": "targetAccount",
				"type": "address"
			},
			{
				"name": "fromTokenNo",
				"type": "string"
			},
			{
				"name": "targetTokenNo",
				"type": "string"
			},
			{
				"name": "formType",
				"type": "uint8"
			},
			{
				"name": "toType",
				"type": "uint8"
			},
			{
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "transfer",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "queryCurrentAccounts",
		"outputs": [
			{
				"name": "",
				"type": "address[]"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "sysAccount",
				"type": "address"
			}
		],
		"name": "setSysAccount",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "_operationAccount",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "",
				"type": "address[]"
			},
			{
				"indexed": false,
				"name": "",
				"type": "address[]"
			},
			{
				"indexed": false,
				"name": "",
				"type": "uint256[]"
			}
		],
		"name": "event_publish",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "",
				"type": "uint256"
			}
		],
		"name": "event_transfer",
		"type": "event"
	}
]`

var testAbi string = `[
	{
		"constant": false,
		"inputs": [
			{
				"name": "newadds",
				"type": "address[]"
			}
		],
		"name": "test1",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "newadds",
				"type": "address[2]"
			}
		],
		"name": "test2",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "newnums",
				"type": "uint256[]"
			}
		],
		"name": "test3",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "newnums",
				"type": "uint256[2]"
			}
		],
		"name": "test4",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"name": "adds",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"name": "nums",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	}
]`

var controlBytecode string = "608060405234801561001057600080fd5b506002805560008054600160a060020a031916331790556102d6806100366000396000f3006080604052600436106100825763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663092a5cce81146100875780630fb3844c146100b057806314b599f6146100d75780635a0bc4d3146100ec578063b2bdfa7b1461011a578063c3bf930e1461014b578063df93ad1314610160575b600080fd5b34801561009357600080fd5b5061009c610180565b604080519115158252519081900360200190f35b3480156100bc57600080fd5b506100c56101a0565b60408051918252519081900360200190f35b3480156100e357600080fd5b5061009c6101a6565b3480156100f857600080fd5b5061009c6024600480358281019291013590600160a060020a039035166101c6565b34801561012657600080fd5b5061012f61024b565b60408051600160a060020a039092168252519081900360200190f35b34801561015757600080fd5b506100c561025a565b34801561016c57600080fd5b5061012f6004803560248101910135610260565b60008054600160a060020a0316331461019857600080fd5b600360025590565b60025481565b60008054600160a060020a031633146101be57600080fd5b600160025590565b60008054600160a060020a031633146101de57600080fd5b60028054146101ec57600080fd5b8160018585604051808383808284379091019485525050604051928390036020019092208054600160a060020a039490941673ffffffffffffffffffffffffffffffffffffffff19909416939093179092555060019150509392505050565b600054600160a060020a031681565b60025490565b6000600280541115151561027357600080fd5b600183836040518083838082843790910194855250506040519283900360200190922054600160a060020a031692505050929150505600a165627a7a72305820c6f490266f6861f04679ab4d019d9d39820ef4c982ceb435657be33a4039f3480029"
var tokenBytecode string = "608060405234801561001057600080fd5b5060008054600160a060020a0319163317905560026003556127aa806100376000396000f3006080604052600436106101695763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663092a5cce811461016e57806314b599f614610197578063150ea47e146101ac5780631889523c146101cd57806321f50b77146101f857806327d92abf1461022457806339dce3dd146103e7578063446614e514610408578063549e52021461043057806354cd903c1461045c57806363aa002f146104905780636b408ebd146104d457806370d38acd146104f55780637622c6501461051d57806381f6658c146105495780639c16d24814610583578063a16e5dc2146105e8578063a2fa231c14610608578063a657399814610628578063b2578b1114610654578063b2bdfa7b146106ed578063b592770614610702578063bed2952e14610735578063cbbd7bed1461077c578063e3174a3e146107a0578063f0b91ee1146107cc578063f5885e3a146107ed578063f7a9039414610824575b600080fd5b34801561017a57600080fd5b50610183610839565b604080519115158252519081900360200190f35b3480156101a357600080fd5b50610183610858565b3480156101b857600080fd5b50610183600160a060020a0360043516610878565b3480156101d957600080fd5b50610183602460048035828101929101359060ff9035166044356108da565b34801561020457600080fd5b506101836024600480358281019290820135918135918201910135610af4565b34801561023057600080fd5b506102446004803560248101910135610b9f565b6040518080602001806020018c600160a060020a0316600160a060020a031681526020018b600160a060020a0316600160a060020a031681526020018a81526020018981526020018881526020018781526020018060200186815260200185815260200184810384528f818151815260200191508051906020019080838360005b838110156102dd5781810151838201526020016102c5565b50505050905090810190601f16801561030a5780820380516001836020036101000a031916815260200191505b5084810383528e818151815260200191508051906020019080838360005b83811015610340578181015183820152602001610328565b50505050905090810190601f16801561036d5780820380516001836020036101000a031916815260200191505b50848103825287518152875160209182019189019080838360005b838110156103a0578181015183820152602001610388565b50505050905090810190601f1680156103cd5780820380516001836020036101000a031916815260200191505b509e50505050505050505050505050505060405180910390f35b3480156103f357600080fd5b50610183600160a060020a0360043516610e86565b34801561041457600080fd5b5061042e600160a060020a03600435166024351515610ee3565b005b34801561043c57600080fd5b506101836024600480358281019290820135918135918201910135610f25565b34801561046857600080fd5b50610183602460048035828101929101359060ff903516600160a060020a03604435166110b5565b34801561049c57600080fd5b506104b8602460048035828101929101359060ff903516611214565b60408051600160a060020a039092168252519081900360200190f35b3480156104e057600080fd5b50610183600160a060020a036004351661133e565b34801561050157600080fd5b50610183602460048035828101929101359060ff903516611353565b34801561052957600080fd5b506101836024600480358281019290820135918135918201910135611436565b34801561055557600080fd5b50610571602460048035828101929101359060ff90351661156d565b60408051918252519081900360200190f35b34801561058f57600080fd5b5061059861176c565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156105d45781810151838201526020016105bc565b505050509050019250505060405180910390f35b3480156105f457600080fd5b50610183600480356024810191013561181e565b34801561061457600080fd5b506104b86004803560248101910135611880565b34801561063457600080fd5b5061018360246004803582810192908201359181359182019101356118eb565b34801561066057600080fd5b506106786024600480358281019291013590356119fa565b6040805160208082528351818301528351919283929083019185019080838360005b838110156106b257818101518382015260200161069a565b50505050905090810190601f1680156106df5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156106f957600080fd5b506104b8611c82565b34801561070e57600080fd5b5061018360246004803582810192908201359160ff82351691604435908101910135611c91565b34801561074157600080fd5b506101836024600480358281019290820135918135918201910135600160a060020a036044358116906064351660843560a43560c435611e0c565b34801561078857600080fd5b50610678602460048035828101929101359035612017565b3480156107ac57600080fd5b506101836024600480358281019290820135918135918201910135612291565b3480156107d857600080fd5b50610183600160a060020a0360043516612459565b3480156107f957600080fd5b50610183602460048035828101929082013591813591604435916064359182019101356084356124a2565b34801561083057600080fd5b506104b86125d1565b60008054600160a060020a0316331461085157600080fd5b6003805590565b60008054600160a060020a0316331461087057600080fd5b600160035590565b600154600090600160a060020a0316331461089257600080fd5b600154600160a060020a031633146108a957600080fd5b6001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03939093169290921790915590565b3360009081526004602052604081205460ff1615156001146108fb57600080fd5b60035460021461090a57600080fd5b60058585604051808383808284379091019485525050604051928390036020019092205460ff1615159150610940905057600080fd5b60058585604051808383808284379091019485525050604051928390036020019092205460ff16159150610ae89050578260ff16600514156109ae57816005868660405180838380828437820191505092505050908152602001604051809103902060050181905550610ae0565b8260ff16600614156109ec57816005868660405180838380828437820191505092505050908152602001604051809103902060060181905550610ae0565b8260ff1660071415610a2a57816005868660405180838380828437820191505092505050908152602001604051809103902060070181905550610ae0565b8260ff1660081415610a6857816005868660405180838380828437820191505092505050908152602001604051809103902060080181905550610ae0565b8260ff16600a1415610aa6578160058686604051808383808284378201915050925050509081526020016040518091039020600a0181905550610ae0565b8260ff16600b1415610ae0578160058686604051808383808284378201915050925050509081526020016040518091039020600b01819055505b506001610aec565b5060005b949350505050565b3360009081526004602052604081205460ff161515600114610b1557600080fd5b610b96610b5184848080601f016020809104026020016040519081016040528093929190818152602001838380828437506125e0945050505050565b610b8a87878080601f016020809104026020016040519081016040528093929190818152602001838380828437506125e0945050505050565b9063ffffffff61260616565b95945050505050565b33600090815260046020526040812054606091829181908190819081908190879082908190819060ff161515600114610bd757600080fd5b60035460021015610be757600080fd5b60058e8e604051808383808284379091019485525050604051928390036020019092205460ff16159150610e7590505760058e8e604051808383808284378201915050925050509081526020016040518091039020905080600101816002018260030160009054906101000a9004600160a060020a03168360040160009054906101000a9004600160a060020a031684600501548560060154866007015487600801548860090189600a01548a600b01548a8054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610d2d5780601f10610d0257610100808354040283529160200191610d2d565b820191906000526020600020905b815481529060010190602001808311610d1057829003601f168201915b50505050509a50898054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610dc95780601f10610d9e57610100808354040283529160200191610dc9565b820191906000526020600020905b815481529060010190602001808311610dac57829003601f168201915b5050865460408051602060026001851615610100026000190190941693909304601f8101849004840282018401909252818152959f5088945092508401905082828015610e575780601f10610e2c57610100808354040283529160200191610e57565b820191906000526020600020905b815481529060010190602001808311610e3a57829003601f168201915b505050505092509b509b509b509b509b509b509b509b509b509b509b505b509295989b509295989b9093969950565b60008054600160a060020a03163314610e9e57600080fd5b600160a060020a0382161515610eb357600080fd5b5060018054600160a060020a03831673ffffffffffffffffffffffffffffffffffffffff19909116178155919050565b600054600160a060020a03163314610efa57600080fd5b600160a060020a03919091166000908152600460205260409020805460ff1916911515919091179055565b3360009081526004602052604081205460ff161515600114610f4657600080fd5b60035460021015610f5657600080fd5b60058585604051808383808284379091019485525050604051928390036020019092205460ff1615159150610f8c905057600080fd5b60058585604051808383808284379091019485525050604051928390036020019092205460ff16159150610aec9050576110ae610ff884848080601f016020809104026020016040519081016040528093929190818152602001838380828437506125e0945050505050565b610b8a60058888604051808383808284379190910194855250506040805160209481900385018120600901805460026001821615610100026000190190911604601f810187900487028301870190935282825290949093509091508301828280156110a45780601f10611079576101008083540402835291602001916110a4565b820191906000526020600020905b81548152906001019060200180831161108757829003601f168201915b50505050506125e0565b9050610aec565b3360009081526004602052604081205460ff1615156001146110d657600080fd5b6003546002146110e557600080fd5b60058585604051808383808284379091019485525050604051928390036020019092205460ff161515915061111b905057600080fd5b60058585604051808383808284379091019485525050604051928390036020019092205460ff16159150610ae89050578260ff16600314156111a957816005868660405180838380828437820191505092505050908152602001604051809103902060030160006101000a815481600160a060020a030219169083600160a060020a03160217905550610ae0565b8260ff1660041415610ae0578160058686604051808383808284379091019485525050604051928390036020019092206004018054600160a060020a039490941673ffffffffffffffffffffffffffffffffffffffff19909416939093179092555050506001610aec565b3360009081526004602052604081205460ff16151560011461123557600080fd5b6003546002101561124557600080fd5b60058484604051808383808284379091019485525050604051928390036020019092205460ff161515915061127b905057600080fd5b60058484604051808383808284379091019485525050604051928390036020019092205460ff161591506113379050578160ff16600314156112f157600584846040518083838082843790910194855250506040519283900360200190922060030154600160a060020a03169250611337915050565b8160ff166004141561133757600584846040518083838082843790910194855250506040519283900360200190922060040154600160a060020a03169250611337915050565b9392505050565b60046020526000908152604090205460ff1681565b3360009081526004602052604081205460ff16151560011461137457600080fd5b6003546002101561138457600080fd5b60058484604051808383808284379091019485525050604051928390036020019092205460ff16151591506113ba905057600080fd5b8160ff16600614156113f857426005858560405180838380828437820191505092505050908152602001604051809103902060060154119050611337565b8160ff166008141561133757426005858560405180838380828437820191505092505050908152602001604051809103902060080154119050611337565b3360009081526004602052604081205460ff16151560011461145757600080fd5b6003546002101561146757600080fd5b60058383604051808383808284379091019485525050604051928390036020019092205460ff161515915061149d905057600080fd5b610b966114d986868080601f016020809104026020016040519081016040528093929190818152602001838380828437506125e0945050505050565b610b8a600586866040518083838082843782019150509250505090815260200160405180910390206001018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156110a45780601f10611079576101008083540402835291602001916110a4565b3360009081526004602052604081205460ff16151560011461158e57600080fd5b6003546002101561159e57600080fd5b60058484604051808383808284379091019485525050604051928390036020019092205460ff16151591506115d4905057600080fd5b60058484604051808383808284379091019485525050604051928390036020019092205460ff161591506113379050578160ff16600514156116405760058484604051808383808284378201915050925050509081526020016040518091039020600501549050611337565b8160ff166006141561167c5760058484604051808383808284378201915050925050509081526020016040518091039020600601549050611337565b8160ff16600714156116b85760058484604051808383808284378201915050925050509081526020016040518091039020600701549050611337565b8160ff16600814156116f45760058484604051808383808284378201915050925050509081526020016040518091039020600801549050611337565b8160ff16600a14156117305760058484604051808383808284378201915050925050509081526020016040518091039020600a01549050611337565b8160ff16600b14156113375760058484604051808383808284378201915050925050509081526020016040518091039020600b01549050611337565b604080516003808252608082019092526060918291906020820183803883395050600080548351939450600160a060020a03169284925081106117ab57fe5b600160a060020a03928316602091820290920101526002548251911690829060019081106117d557fe5b600160a060020a03928316602091820290920101526001548251911690829060029081106117ff57fe5b600160a060020a039092166020928302909101909101529050805b5090565b3360009081526004602052604081205460ff16151560011461183f57600080fd5b6003546002101561184f57600080fd5b60058383604051808383808284379091019485525050604051928390036020019092205460ff169250505092915050565b3360009081526004602052604081205460ff1615156001146118a157600080fd5b600354600210156118b157600080fd5b600583836040518083838082843790910194855250506040519283900360200190922060030154600160a060020a03169250505092915050565b3360009081526004602052604081205460ff16151560011461190c57600080fd5b6003546002101561191c57600080fd5b60058585604051808383808284379091019485525050604051928390036020019092205460ff1615159150611952905057600080fd5b60058383604051808383808284379091019485525050604051928390036020019092205460ff1615159150611988905057600080fd5b6005838360405180838380828437909101948552505060405192839003602001832060030154600160a060020a0316926005925088915087908083838082843790910194855250506040519283900360200190922060030154600160a060020a03169290921492505050949350505050565b3360009081526004602052604090205460609060ff161515600114611a1e57600080fd5b60035460021015611a2e57600080fd5b60058484604051808383808284379091019485525050604051928390036020019092205460ff1615159150611a64905057600080fd5b60058484604051808383808284379091019485525050604051928390036020019092205460ff161591506113379050578160011415611b6d57611b66600585856040518083838082843782019150509250505090815260200160405180910390206001018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015611b5c5780601f10611b3157610100808354040283529160200191611b5c565b820191906000526020600020905b815481529060010190602001808311611b3f57829003601f168201915b505050505061261a565b9050611337565b8160021415611bf857611b666005858560405180838380828437919091019485525050604080516020948190038501812060029081018054600181161561010002600019011691909104601f81018790048702830187019093528282529094909350909150830182828015611b5c5780601f10611b3157610100808354040283529160200191611b5c565b816009141561133757610aec60058585604051808383808284379190910194855250506040805160209481900385018120600901805460026001821615610100026000190190911604601f81018790048702830187019093528282529094909350909150830182828015611b5c5780601f10611b3157610100808354040283529160200191611b5c565b600054600160a060020a031681565b3360009081526004602052604081205460ff161515600114611cb257600080fd5b600354600214611cc157600080fd5b60058686604051808383808284379091019485525050604051928390036020019092205460ff1615159150611cf7905057600080fd5b60058686604051808383808284379091019485525050604051928390036020019092205460ff16159150611e009050578360ff1660011415611d70578282600588886040518083838082843782019150509250505090815260200160405180910390206001019190611d6a9291906126d3565b50611df8565b8360ff1660021415611db3578282600588886040518083838082843782019150509250505090815260200160405180910390206002019190611d6a9291906126d3565b8360ff1660091415611df8578282600588886040518083838082843782019150509250505090815260200160405180910390206009019190611df69291906126d3565b505b506001610b96565b50600095945050505050565b3360009081526004602052604081205460ff161515600114611e2d57600080fd5b600354600214611e3c57600080fd5b600160058989604051808383808284379091019485525050604051928390036020018320805494151560ff1990951694909417909355508b91508a906005908b908b908083838082843782019150509250505090815260200160405180910390206001019190611ead9291906126d3565b50878760058a8a6040518083838082843782019150509250505090815260200160405180910390206002019190611ee59291906126d3565b50856005898960405180838380828437820191505092505050908152602001604051809103902060030160006101000a815481600160a060020a030219169083600160a060020a03160217905550846005898960405180838380828437820191505092505050908152602001604051809103902060040160006101000a815481600160a060020a030219169083600160a060020a03160217905550836005898960405180838380828437820191505092505050908152602001604051809103902060050181905550826005898960405180838380828437820191505092505050908152602001604051809103902060060181905550816005898960405180838380828437909101948552505060405192839003602001909220600b01929092555060019b9a5050505050505050505050565b3360009081526004602052604090205460609060ff16151560011461203b57600080fd5b6003546002101561204b57600080fd5b60058484604051808383808284379091019485525050604051928390036020019092205460ff1615159150612081905057600080fd5b60058484604051808383808284379091019485525050604051928390036020019092205460ff16159150611337905057816001141561218257600584846040518083838082843782019150509250505090815260200160405180910390206001018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156121765780601f1061214b57610100808354040283529160200191612176565b820191906000526020600020905b81548152906001019060200180831161215957829003601f168201915b50505050509050611337565b816002141561220a576005848460405180838380828437919091019485525050604080516020948190038501812060029081018054600181161561010002600019011691909104601f810187900487028301870190935282825290949093509091508301828280156121765780601f1061214b57610100808354040283529160200191612176565b81600914156113375760058484604051808383808284379190910194855250506040805160209481900385018120600901805460026001821615610100026000190190911604601f810187900487028301870190935282825290949093509091508301828280156121765780601f1061214b57610100808354040283529160200191612176565b3360009081526004602052604081205460ff1615156001146122b257600080fd5b600354600210156122c257600080fd5b60058585604051808383808284379091019485525050604051928390036020019092205460ff16151591506122f8905057600080fd5b60058383604051808383808284379091019485525050604051928390036020019092205460ff161515915061232e905057600080fd5b610b966123c5600585856040518083838082843782019150509250505090815260200160405180910390206001018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156110a45780601f10611079576101008083540402835291602001916110a4565b610b8a600588886040518083838082843782019150509250505090815260200160405180910390206001018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156110a45780601f10611079576101008083540402835291602001916110a4565b60008054600160a060020a0316331461247157600080fd5b6002805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03939093169290921790915590565b3360009081526004602052604081205460ff1615156001146124c357600080fd5b6003546002146124d257600080fd5b60058888604051808383808284379091019485525050604051928390036020019092205460ff161591506125c6905057856005898960405180838380828437820191505092505050908152602001604051809103902060070181905550846005898960405180838380828437820191505092505050908152602001604051809103902060080181905550838360058a8a60405180838380828437820191505092505050908152602001604051809103902060090191906125939291906126d3565b508160058989604051808383808284378201915050925050509081526020016040518091039020600a0181905550600190505b979650505050505050565b600154600160a060020a031681565b6125e861274d565b50604080518082019091528151815260209182019181019190915290565b6000612612838361261d565b159392505050565b90565b60008060008060008060008060008a6000015197508a600001518a60000151101561264757895197505b8a60200151965089602001519550600094505b878510156126bd578651865190945092508284146126a9576000199150602088101561269357600185896020030160080260020a031991505b508181168382160380156126a9578098506126c5565b60209687019695860195949094019361265a565b89518b510398505b505050505050505092915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106127145782800160ff19823516178555612741565b82800160010185558215612741579182015b82811115612741578235825591602001919060010190612726565b5061181a929150612764565b604080518082019091526000808252602082015290565b61261a91905b8082111561181a576000815560010161276a5600a165627a7a7230582079ae35af5f752ed05d4714f91ff3944c4c72236050b0a0e6932cfa4f0a19de4f0029"
var accountBytecode string = "608060405234801561001057600080fd5b5060008054600160a060020a031916331790556002600155611d4c806100376000396000f3006080604052600436106101485763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630586ef24811461014d57806307f5a1fb146101ca578063092a5cce146101f85780630fb3844c1461022157806314b599f614610248578063150ea47e1461025d57806339dce3dd1461027e578063428437921461029f578063488f646f146102bf57806356da8100146102f05780635a2c9009146103165780636b408ebd146103375780638e0eaa9914610358578063924d1c5b146103945780639c16d248146103af578063a5a8d2fa146103c4578063b2bdfa7b14610428578063c3bf930e1461043d578063d302651f14610452578063da8a97e214610488578063e22d28451461052b578063e9fac80d14610543578063ea37551114610570578063f0b91ee1146105a9578063f3a74977146105ca575b600080fd5b34801561015957600080fd5b5061017a60048035600160a060020a031690602480359081019101356105fd565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156101b657818101518382015260200161019e565b505050509050019250505060405180910390f35b3480156101d657600080fd5b506101f660246004803582810192908201359181359182019101356106d8565b005b34801561020457600080fd5b5061020d6107d6565b604080519115158252519081900360200190f35b34801561022d57600080fd5b506102366107f6565b60408051918252519081900360200190f35b34801561025457600080fd5b5061020d6107fc565b34801561026957600080fd5b5061020d600160a060020a036004351661081b565b34801561028a57600080fd5b5061020d600160a060020a036004351661087f565b3480156102ab57600080fd5b5061017a60048035602481019101356108e0565b3480156102cb57600080fd5b506102d4610993565b60408051600160a060020a039092168252519081900360200190f35b3480156102fc57600080fd5b506101f6600160a060020a036004351660243515156109a2565b34801561032257600080fd5b5061020d600160a060020a03600435166109f3565b34801561034357600080fd5b5061020d600160a060020a0360043516610b07565b34801561036457600080fd5b506101f660048035600160a060020a0316906024803590810191013560443560643560843560a43560c435610b1c565b3480156103a057600080fd5b5061017a600435602435610db8565b3480156103bb57600080fd5b5061017a610ead565b3480156103d057600080fd5b506103fd60048035600160a060020a03169060248035808201929081013591604435908101910135610f5e565b6040805195865260208601949094528484019290925260608401526080830152519081900360a00190f35b34801561043457600080fd5b506102d4611007565b34801561044957600080fd5b50610236611016565b34801561045e57600080fd5b5061020d60048035600160a060020a0316906024803590810191013560443560ff1660643561101d565b34801561049457600080fd5b506104ac600160a060020a0360043516602435611273565b6040518080602001838152602001828103825284818151815260200191508051906020019080838360005b838110156104ef5781810151838201526020016104d7565b50505050905090810190601f16801561051c5780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b34801561053757600080fd5b506102d460043561148d565b34801561054f57600080fd5b5061020d60048035600160a060020a031690602480359081019101356114b5565b34801561057c57600080fd5b506101f660048035600160a060020a031690602480358082019290810135916044359081019101356115c0565b3480156105b557600080fd5b5061020d600160a060020a036004351661169a565b3480156105d657600080fd5b5061023660048035600160a060020a0316906024803590810191013560443560ff166116e5565b3360009081526009602052604090205460609060ff16151560011461062157600080fd5b6001546002101561063157600080fd5b6006600085600160a060020a0316600160a060020a03168152602001908152602001600020600101838360405180838380828437909101948552505060408051938490036020908101852080548083028701830190935282865293509091508301828280156106c957602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116106ab575b505050505090505b9392505050565b600354606090600090600160a060020a031633146106f557600080fd5b60015460021461070457600080fd5b600a83111561071257600080fd5b60408051848152602080860282010190915283801561073b578160200160208202803883390190505b509150600090505b828110156107945783838281811061075757fe5b90506020020135600160a060020a0316828281518110151561077557fe5b600160a060020a03909216602092830290910190910152600101610743565b816007878760405180838380828437820191505092505050908152602001604051809103902090805190602001906107cd929190611b9d565b50505050505050565b60008054600160a060020a031633146107ee57600080fd5b600360015590565b60015481565b60008054600160a060020a0316331461081457600080fd5b6001805590565b600254600090600160a060020a0316331461083557600080fd5b60015460021461084457600080fd5b600254600160a060020a0316331461085b57600080fd5b60028054600160a060020a031916600160a060020a03939093169290921790915590565b60008054600160a060020a0316331461089757600080fd5b6001546002146108a657600080fd5b600160a060020a03821615156108bb57600080fd5b5060028054600160a060020a038316600160a060020a03199091161790556001919050565b3360009081526009602052604090205460609060ff16151560011461090457600080fd5b6001546002101561091457600080fd5b60078383604051808383808284379091019485525050604080519384900360209081018520805480830287018301909352828652935090915083018282801561098657602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610968575b5050505050905092915050565b600254600160a060020a031681565b600054600160a060020a031633146109b957600080fd5b6001546002146109c857600080fd5b600160a060020a03919091166000908152600960205260409020805460ff1916911515919091179055565b60008054600160a060020a03163314610a0b57600080fd5b600a8054600160a060020a031916600160a060020a038481169190911791829055604080517fdf93ad13000000000000000000000000000000000000000000000000000000008152602060048201819052601160248301527f4173736574546f6b656e4164647265737300000000000000000000000000000060448301529151939092169263df93ad139260648082019392918290030181600087803b158015610ab457600080fd5b505af1158015610ac8573d6000803e3d6000fd5b505050506040513d6020811015610ade57600080fd5b505160048054600160a060020a031916600160a060020a03909216919091179055506001919050565b60096020526000908152604090205460ff1681565b3360009081526009602052604090205460ff161515600114610b3d57600080fd5b600154600214610b4c57600080fd5b610b598585858585611908565b610b6288611955565b600160a060020a03881660008181526005602052604090819020805474ffffffffffffffffffffffffffffffffffffffff00191661010090930292909217825551889188916002909101908390839080838380828437909101948552505060405192839003602001909220610bdd9490939092509050611bfe565b5084600560008a600160a060020a0316600160a060020a0316815260200190815260200160002060020188886040518083838082843782019150509250505090815260200160405180910390206001018190555083600560008a600160a060020a0316600160a060020a031681526020019081526020016000206002018888604051808383808284379091019485525050604080519384900360209081018520600290810196909655600160a060020a038e166000908152600590915220879401928b92508a91508083838082843782019150509250505090815260200160405180910390206003018190555081600560008a600160a060020a0316600160a060020a0316815260200190815260200160002060020188886040518083838082843782019150509250505090815260200160405180910390206004018190555080600560008a600160a060020a0316600160a060020a031681526020019081526020016000206002018888604051808383808284379190910194855250506040805160209481900385019020600590810195909555600160a060020a038d166000908152948452842060019081018054918201808255908652939094209293610dac930191508a905089611bfe565b50505050505050505050565b3360009081526009602052604081205460609190819083908290819060ff161515600114610de557600080fd5b60015460021015610df557600080fd5b600854945086851115610e085786610e0a565b845b935083604051908082528060200260200182016040528015610e36578160200160208202803883390190505b509250600091505b83821015610ea15760088054898401908110610e5657fe5b6000918252602090912001548351600160a060020a0390911691508190849084908110610e7f57fe5b600160a060020a03909216602092830290910190910152600190910190610e3e565b50909695505050505050565b604080516003808252608082019092526060918291906020820183803883395050600080548351939450600160a060020a0316928492508110610eec57fe5b600160a060020a0392831660209182029092010152600354825191169082906001908110610f1657fe5b600160a060020a03928316602091820290920101526002805483519216918391908110610f3f57fe5b600160a060020a039092166020928302909101909101529050805b5090565b336000908152600960205260408120548190819081908190819060ff161515600114610f8957600080fd5b600560008c600160a060020a0316600160a060020a03168152602001908152602001600020600201888860405180838380828437820191505092505050908152602001604051809103902090508060010154816002015482600301548360040154846005015495509550955095509550509550955095509550959050565b600054600160a060020a031681565b6001545b90565b3360009081526009602052604081205460ff16151560011461103e57600080fd5b60015460021461104d57600080fd5b600160a060020a03861660009081526005602052604090205460ff1615611266578260ff16600114156110d257816005600088600160a060020a0316600160a060020a0316815260200190815260200160002060020186866040518083838082843782019150509250505090815260200160405180910390206001018190555061125e565b8260ff166002141561113657816005600088600160a060020a0316600160a060020a0316815260200190815260200160002060020186866040518083838082843782019150509250505090815260200160405180910390206002018190555061125e565b8260ff166003141561119a57816005600088600160a060020a0316600160a060020a0316815260200190815260200160002060020186866040518083838082843782019150509250505090815260200160405180910390206003018190555061125e565b8260ff16600414156111fe57816005600088600160a060020a0316600160a060020a0316815260200190815260200160002060020186866040518083838082843782019150509250505090815260200160405180910390206004018190555061125e565b8260ff166005141561125e57816005600088600160a060020a0316600160a060020a031681526020019081526020016000206002018686604051808383808284378201915050925050509081526020016040518091039020600501819055505b50600161126a565b5060005b95945050505050565b33600090815260096020526040812054606091908290829082908290819060ff1615156001146112a257600080fd5b600154600210156112b257600080fd5b600160a060020a038916600090815260056020908152604080832060010180548251818502810185019093528083529193909284015b828210156113935760008481526020908190208301805460408051601f600260001961010060018716150201909416939093049283018590048502810185019091528181529283018282801561137f5780601f106113545761010080835404028352916020019161137f565b820191906000526020600020905b81548152906001019060200180831161136257829003601f168201915b5050505050815260200190600101906112e8565b50505050945084519350600088101580156113ad57508388105b1561147e57600588850310156113c5578784036113c8565b60055b91508790505b81880181101561147e578781141561141e5761141861140386838151811015156113f457fe5b90602001906020020151611a0d565b61140c85611a0d565b9063ffffffff611a3316565b50611476565b61145f6114036040805190810160405280600181526020017f2600000000000000000000000000000000000000000000000000000000000000815250611a0d565b5061147461140386838151811015156113f457fe5b505b6001016113ce565b50909791965090945050505050565b600880548290811061149b57fe5b600091825260209091200154600160a060020a0316905081565b60006002600154111515156114c957600080fd5b6115aa6115a56005600087600160a060020a0316600160a060020a031681526020019081526020016000206002018585604051808383808284379190910194855250506040805160209481900385018120805460026001821615610100026000190190911604601f8101879004870283018701909352828252909490935090915083018282801561159b5780601f106115705761010080835404028352916020019161159b565b820191906000526020600020905b81548152906001019060200180831161157e57829003601f168201915b5050505050611a0d565b611aaa565b15156115b8575060006106d1565b5060016106d1565b3360009081526009602052604081205460ff1615156001146115e157600080fd5b6001546002146115f057600080fd5b600a8211156115fe57600080fd5b50600160a060020a038516600081815260066020526040908190208054600160a060020a031916909217825551839083906001840190889088908083838082843782019150509250505090815260200160405180910390209190611663929190611c78565b50600160a060020a03958616600090815260066020526040902090548154600160a060020a03191696169590951790945550505050565b60008054600160a060020a031633146116b257600080fd5b6001546002146116c157600080fd5b60038054600160a060020a031916600160a060020a03939093169290921790915590565b60006002600154111515156116f957600080fd5b600160a060020a03851660009081526005602052604090205460ff1615611900578160ff166001141561177c576005600086600160a060020a0316600160a060020a031681526020019081526020016000206002018484604051808383808284378201915050925050509081526020016040518091039020600101549050611900565b8160ff16600214156117de576005600086600160a060020a0316600160a060020a031681526020019081526020016000206002018484604051808383808284378201915050925050509081526020016040518091039020600201549050611900565b8160ff1660031415611840576005600086600160a060020a0316600160a060020a031681526020019081526020016000206002018484604051808383808284378201915050925050509081526020016040518091039020600301549050611900565b8160ff16600414156118a2576005600086600160a060020a0316600160a060020a031681526020019081526020016000206002018484604051808383808284378201915050925050509081526020016040518091039020600401549050611900565b8160ff1660051415611900576005600086600160a060020a0316600160a060020a0316815260200190815260200160002060020184846040518083838082843782019150509250505090815260200160405180910390206005015490505b949350505050565b600085101561191657600080fd5b600084101561192457600080fd5b600083101561193257600080fd5b600082101561194057600080fd5b600081101561194e57600080fd5b5050505050565b6001546002101561196557600080fd5b600160a060020a03811660009081526005602052604090205460ff161515611a0a576008805460018082019092557ff3f7a9fe364faab93b216da50a3214154f22a0a2b415b23a84c8169e8b636ee3018054600160a060020a031916600160a060020a0384169081179091556000818152600560205260409020805460ff191690921774ffffffffffffffffffffffffffffffffffffffff0019166101009091021790555b50565b611a15611ccb565b50604080518082019091528151815260209182019181019190915290565b606080600083600001518560000151016040519080825280601f01601f191660200182016040528015611a70578160200160208202803883390190505b509150602082019050611a8c8186602001518760000151611b59565b845160208501518551611aa29284019190611b59565b509392505050565b60208101518151600091601e198082019290910101825b81831015611b515750815160ff166080811015611ae357600183019250611b46565b60e08160ff161015611afa57600283019250611b46565b60f08160ff161015611b1157600383019250611b46565b60f88160ff161015611b2857600483019250611b46565b60fc8160ff161015611b3f57600583019250611b46565b6006830192505b600190930192611ac1565b505050919050565b60005b60208210611b7e578251845260209384019390920191601f1990910190611b5c565b50905182516020929092036101000a6000190180199091169116179052565b828054828255906000526020600020908101928215611bf2579160200282015b82811115611bf25782518254600160a060020a031916600160a060020a03909116178255602090920191600190910190611bbd565b50610f5a929150611ce2565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10611c3f5782800160ff19823516178555611c6c565b82800160010185558215611c6c579182015b82811115611c6c578235825591602001919060010190611c51565b50610f5a929150611d06565b828054828255906000526020600020908101928215611bf2579160200282015b82811115611bf2578154600160a060020a031916600160a060020a03843516178255602090920191600190910190611c98565b604080518082019091526000808252602082015290565b61101a91905b80821115610f5a578054600160a060020a0319168155600101611ce8565b61101a91905b80821115610f5a5760008155600101611d0c5600a165627a7a72305820b09cdd1a1d618ab654d00ac621f0de6b29337fffb27aed61e1710725cc8c28170029"
var transBytecode string = "608060405234801561001057600080fd5b506000805433600160a060020a0319918216811783556001805490921617905561291f90819061004090396000f30060806040526004361061008a5763ffffffff60e060020a600035041663050d529c811461008f578063092a5cce146100da57806314b599f614610103578063150ea47e1461011857806339dce3dd146101395780635a2c90091461015a5780637d1c22ff1461017b5780639c16d24814610234578063f0b91ee114610299578063f7a90394146102ba575b600080fd5b34801561009b57600080fd5b506100d8600480359060248035600160a060020a03169160443580830192908201359160643591820191013560843560a43560c43560e4356102eb565b005b3480156100e657600080fd5b506100ef6106d5565b604080519115158252519081900360200190f35b34801561010f57600080fd5b506100ef6106f4565b34801561012457600080fd5b506100ef600160a060020a0360043516610714565b34801561014557600080fd5b506100ef600160a060020a0360043516610776565b34801561016657600080fd5b506100ef600160a060020a03600435166107d3565b34801561018757600080fd5b50604080516020600460443581810135601f81018490048402850184019095528484526100d89482359460248035600160a060020a03169536959460649492019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a9998810197919650918201945092508291508401838280828437509497505060ff853581169650602086013516946040013593506109e492505050565b34801561024057600080fd5b506102496110a8565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561028557818101518382015260200161026d565b505050509050019250505060405180910390f35b3480156102a557600080fd5b506100ef600160a060020a0360043516611157565b3480156102c657600080fd5b506102cf6111a0565b60408051600160a060020a039092168252519081900360200190f35b600254600160a060020a0316331461030257600080fd5b60035460021461031157600080fd5b600061031e8460006111af565b10158015610337575060006103348360006111af565b10155b8015610354575066038d7ea4c680006103518460006111af565b11155b8015610371575066038d7ea4c6800061036e8360006111af565b11155b151561037c57600080fd5b80600114156104fe576006546040517fa16e5dc200000000000000000000000000000000000000000000000000000000815260206004820190815260248201889052600160a060020a039092169163a16e5dc29189918991819060440184848082843782019150509350505050602060405180830381600087803b15801561040357600080fd5b505af1158015610417573d6000803e3d6000fd5b505050506040513d602081101561042d57600080fd5b50511561043957600080fd5b600061044584846111c2565b101561045057600080fd5b6104ba8989898080601f0160208091040260200160405190810160405280939291908181526020018383808284375050604080516020601f8f018190048102820181019092528d815294508d93508c92508291508401838280828437508c94506111cf9350505050565b506104f98987878080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050508585611471565b6106c9565b6006546040517fa16e5dc200000000000000000000000000000000000000000000000000000000815260206004820190815260248201889052600160a060020a039092169163a16e5dc29189918991819060440184848082843782019150509350505050602060405180830381600087803b15801561057c57600080fd5b505af1158015610590573d6000803e3d6000fd5b505050506040513d60208110156105a657600080fd5b505115156105b357600080fd5b600654604080517f7622c65000000000000000000000000000000000000000000000000000000000815260048101918252604481018a9052600160a060020a0390921691637622c650918b918b918b918b9190819060248101906064018787808284379091018481038352858152602001905085858082843782019150509650505050505050602060405180830381600087803b15801561065357600080fd5b505af1158015610667573d6000803e3d6000fd5b505050506040513d602081101561067d57600080fd5b5051151561068a57600080fd5b6106c98a8a88888080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050508486611b40565b50505050505050505050565b60008054600160a060020a031633146106ed57600080fd5b6003805590565b60008054600160a060020a0316331461070c57600080fd5b600160035590565b600154600090600160a060020a0316331461072e57600080fd5b600154600160a060020a0316331461074557600080fd5b6001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03939093169290921790915590565b60008054600160a060020a0316331461078e57600080fd5b600160a060020a03821615156107a357600080fd5b5060018054600160a060020a03831673ffffffffffffffffffffffffffffffffffffffff19909116178155919050565b60008054600160a060020a031633146107eb57600080fd5b6005805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a038481169190911791829055604080517fdf93ad13000000000000000000000000000000000000000000000000000000008152602060048201819052601660248301527f4163636f756e74436f6e7472616374416464726573730000000000000000000060448301529151939092169263df93ad139260648082019392918290030181600087803b1580156108a157600080fd5b505af11580156108b5573d6000803e3d6000fd5b505050506040513d60208110156108cb57600080fd5b50516004805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03928316178155600554604080517fdf93ad130000000000000000000000000000000000000000000000000000000081526020938101849052601160248201527f4173736574546f6b656e4164647265737300000000000000000000000000000060448201529051919093169263df93ad139260648083019391928290030181600087803b15801561098057600080fd5b505af1158015610994573d6000803e3d6000fd5b505050506040513d60208110156109aa57600080fd5b50516006805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0390921691909117905550506002600355600190565b6000806109f488888888886121cd565b60008310158015610a0c575066038d7ea4c680008311155b1515610a1757600080fd5b6004805460405160e060020a63f3a749770281523392810183815260ff891660448301526060602483019081528b5160648401528b51600160a060020a039094169463f3a749779490938d938c939092909160840190602086019080838360005b83811015610a90578181015183820152602001610a78565b50505050905090810190601f168015610abd5780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b158015610ade57600080fd5b505af1158015610af2573d6000803e3d6000fd5b505050506040513d6020811015610b0857600080fd5b5051915082821015610b1957600080fd5b60048054604080517fe9fac80d000000000000000000000000000000000000000000000000000000008152600160a060020a038c8116948201948552602482019283528a5160448301528a5193169363e9fac80d938d938c939091606490910190602085019080838360005b83811015610b9d578181015183820152602001610b85565b50505050905090810190601f168015610bca5780820380516001836020036101000a031916815260200191505b509350505050602060405180830381600087803b158015610bea57600080fd5b505af1158015610bfe573d6000803e3d6000fd5b505050506040513d6020811015610c1457600080fd5b50511515610d3a57600460009054906101000a9004600160a060020a0316600160a060020a0316638e0eaa99898860008060008060006040518863ffffffff1660e060020a0281526004018088600160a060020a0316600160a060020a0316815260200180602001878152602001868152602001858152602001848152602001838152602001828103825288818151815260200191508051906020019080838360005b83811015610ccf578181015183820152602001610cb7565b50505050905090810190601f168015610cfc5780820380516001836020036101000a031916815260200191505b5098505050505050505050600060405180830381600087803b158015610d2157600080fd5b505af1158015610d35573d6000803e3d6000fd5b505050505b6004805460405160e060020a63f3a74977028152600160a060020a038b811693820193845260ff881660448301526060602483019081528a5160648401528a51919093169363f3a74977938d938c938b93929160840190602086019080838360005b83811015610db4578181015183820152602001610d9c565b50505050905090810190601f168015610de15780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b158015610e0257600080fd5b505af1158015610e16573d6000803e3d6000fd5b505050506040513d6020811015610e2c57600080fd5b5051600454909150600160a060020a031663d302651f338988610e4f87896111c2565b6040518563ffffffff1660e060020a0281526004018085600160a060020a0316600160a060020a03168152602001806020018460ff1660ff168152602001838152602001828103825285818151815260200191508051906020019080838360005b83811015610ec8578181015183820152602001610eb0565b50505050905090810190601f168015610ef55780820380516001836020036101000a031916815260200191505b5095505050505050602060405180830381600087803b158015610f1757600080fd5b505af1158015610f2b573d6000803e3d6000fd5b505050506040513d6020811015610f4157600080fd5b5050600454600160a060020a031663d302651f898887610f6186896111af565b6040518563ffffffff1660e060020a0281526004018085600160a060020a0316600160a060020a03168152602001806020018460ff1660ff168152602001838152602001828103825285818151815260200191508051906020019080838360005b83811015610fda578181015183820152602001610fc2565b50505050905090810190601f1680156110075780820380516001836020036101000a031916815260200191505b5095505050505050602060405180830381600087803b15801561102957600080fd5b505af115801561103d573d6000803e3d6000fd5b505050506040513d602081101561105357600080fd5b505060408051338152600160a060020a038a16602082015280820185905290517fe03f3c77300d5d7f8d0abcb05df541a43155707e2a53eda144b929631ef939e59181900360600190a1505050505050505050565b604080516003808252608082019092526060918291906020820183803883395050600080548351939450600160a060020a03169284925081106110e757fe5b600160a060020a039283166020918202909201015260025482519116908290600190811061111157fe5b600160a060020a039283166020918202909201015260015482519116908290600290811061113b57fe5b600160a060020a03909216602092830290910190910152905090565b60008054600160a060020a0316331461116f57600080fd5b6002805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03939093169290921790915590565b600154600160a060020a031681565b808201828110156111bc57fe5b92915050565b808203828111156111bc57fe5b6006546040517fbed2952e000000000000000000000000000000000000000000000000000000008152600160a060020a0386811660448301526000606483018190526084830181905260a4830181905260c4830185905260e060048401908152875160e485015287519194929092169263bed2952e92889288928b928892839283928c92909182916024820191610104019060208c01908083838a5b8381101561128357818101518382015260200161126b565b50505050905090810190601f1680156112b05780820380516001836020036101000a031916815260200191505b5083810382528951815289516020918201918b019080838360005b838110156112e35781810151838201526020016112cb565b50505050905090810190601f1680156113105780820380516001836020036101000a031916815260200191505b509950505050505050505050602060405180830381600087803b15801561133657600080fd5b505af115801561134a573d6000803e3d6000fd5b505050506040513d602081101561136057600080fd5b50506006546040517f1889523c000000000000000000000000000000000000000000000000000000008152600a602482018190524260448301819052606060048401908152875160648501528751600160a060020a0390951694631889523c948994939291829160840190602087019080838360005b838110156113ee5781810151838201526020016113d6565b50505050905090810190601f16801561141b5780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b15801561143c57600080fd5b505af1158015611450573d6000803e3d6000fd5b505050506040513d602081101561146657600080fd5b509095945050505050565b60408051600280825260608281019093528291829181602001602082028038833950506040805160028082526060820183529396509291506020830190803883395050604080516002808252606082018352939550929150602083019080388339019050509050600460009054906101000a9004600160a060020a0316600160a060020a0316638e0eaa99888860008060008060006040518863ffffffff1660e060020a0281526004018088600160a060020a0316600160a060020a0316815260200180602001878152602001868152602001858152602001848152602001838152602001828103825288818151815260200191508051906020019080838360005b8381101561158b578181015183820152602001611573565b50505050905090810190601f1680156115b85780820380516001836020036101000a031916815260200191505b5098505050505050505050600060405180830381600087803b1580156115dd57600080fd5b505af11580156115f1573d6000803e3d6000fd5b5050600480546040517fd302651f000000000000000000000000000000000000000000000000000000008152600160a060020a038c8116938201938452600160448301819052606483018c90526080602484019081528d5160848501528d5192909416965063d302651f95508d948d9491938d9391929160a490910190602087019080838360005b83811015611691578181015183820152602001611679565b50505050905090810190601f1680156116be5780820380516001836020036101000a031916815260200191505b5095505050505050602060405180830381600087803b1580156116e057600080fd5b505af11580156116f4573d6000803e3d6000fd5b505050506040513d602081101561170a57600080fd5b5050825133908490600090811061171d57fe5b600160a060020a03909216602092830290910190910152815187908390600090811061174557fe5b600160a060020a03909216602092830290910190910152805185908290600090811061176d57fe5b60209081029091010152600454600160a060020a031663d302651f888860016117968a8a6111c2565b6040518563ffffffff1660e060020a0281526004018085600160a060020a0316600160a060020a03168152602001806020018460ff168152602001838152602001828103825285818151815260200191508051906020019080838360005b8381101561180c5781810151838201526020016117f4565b50505050905090810190601f1680156118395780820380516001836020036101000a031916815260200191505b5095505050505050602060405180830381600087803b15801561185b57600080fd5b505af115801561186f573d6000803e3d6000fd5b505050506040513d602081101561188557600080fd5b5050600480546040517fd302651f000000000000000000000000000000000000000000000000000000008152600160a060020a038a8116938201938452600260448301819052606483018990526080602484019081528b5160848501528b51929094169463d302651f948d948d948c93929160a40190602087019080838360005b8381101561191e578181015183820152602001611906565b50505050905090810190601f16801561194b5780820380516001836020036101000a031916815260200191505b5095505050505050602060405180830381600087803b15801561196d57600080fd5b505af1158015611981573d6000803e3d6000fd5b505050506040513d602081101561199757600080fd5b50508151829060009081106119a857fe5b906020019060200201518360018151811015156119c157fe5b600160a060020a039092166020928302909101909101528151829060009081106119e757fe5b90602001906020020151826001815181101515611a0057fe5b600160a060020a039092166020928302909101909101528051849082906001908110611a2857fe5b90602001906020020181815250507fb2e1acee2fad7a352b9b70b4b0e6d2f18623c5110f671036a62b2819499b015f83838360405180806020018060200180602001848103845287818151815260200191508051906020019060200280838360005b83811015611aa2578181015183820152602001611a8a565b50505050905001848103835286818151815260200191508051906020019060200280838360005b83811015611ae1578181015183820152602001611ac9565b50505050905001848103825285818151815260200191508051906020019060200280838360005b83811015611b20578181015183820152602001611b08565b50505050905001965050505050505060405180910390a150505050505050565b6004805460405160e060020a63f3a74977028152600160a060020a038088169382019384526001604483018190526060602484019081528851606485015288518a96600096879695169463f3a749779489948d94919392909160849091019060208601908083838d5b83811015611bc1578181015183820152602001611ba9565b50505050905090810190601f168015611bee5780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b158015611c0f57600080fd5b505af1158015611c23573d6000803e3d6000fd5b505050506040513d6020811015611c3957600080fd5b50516004805460405160e060020a63f3a74977028152600160a060020a038781169382019384526002604483018190526060602484019081528c5160648501528c51969850919093169463f3a749779489948d9490939192608490910190602086019080838360005b83811015611cba578181015183820152602001611ca2565b50505050905090810190601f168015611ce75780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b158015611d0857600080fd5b505af1158015611d1c573d6000803e3d6000fd5b505050506040513d6020811015611d3257600080fd5b505190506002851415611f775783821015611d4c57600080fd5b600454600160a060020a031663d302651f84886001611d6b878a6111c2565b6040518563ffffffff1660e060020a0281526004018085600160a060020a0316600160a060020a03168152602001806020018460ff168152602001838152602001828103825285818151815260200191508051906020019080838360005b83811015611de1578181015183820152602001611dc9565b50505050905090810190601f168015611e0e5780820380516001836020036101000a031916815260200191505b5095505050505050602060405180830381600087803b158015611e3057600080fd5b505af1158015611e44573d6000803e3d6000fd5b505050506040513d6020811015611e5a57600080fd5b5050600454600160a060020a031663d302651f84886002611e7b868a6111af565b6040518563ffffffff1660e060020a0281526004018085600160a060020a0316600160a060020a03168152602001806020018460ff168152602001838152602001828103825285818151815260200191508051906020019080838360005b83811015611ef1578181015183820152602001611ed9565b50505050905090810190601f168015611f1e5780820380516001836020036101000a031916815260200191505b5095505050505050602060405180830381600087803b158015611f4057600080fd5b505af1158015611f54573d6000803e3d6000fd5b505050506040513d6020811015611f6a57600080fd5b50611f7790508385612717565b84600314156121b85783811015611f8d57600080fd5b600454600160a060020a031663d302651f84886001611fac878a6111af565b6040518563ffffffff1660e060020a0281526004018085600160a060020a0316600160a060020a03168152602001806020018460ff168152602001838152602001828103825285818151815260200191508051906020019080838360005b8381101561202257818101518382015260200161200a565b50505050905090810190601f16801561204f5780820380516001836020036101000a031916815260200191505b5095505050505050602060405180830381600087803b15801561207157600080fd5b505af1158015612085573d6000803e3d6000fd5b505050506040513d602081101561209b57600080fd5b5050600454600160a060020a031663d302651f848860026120bc868a6111c2565b6040518563ffffffff1660e060020a0281526004018085600160a060020a0316600160a060020a03168152602001806020018460ff168152602001838152602001828103825285818151815260200191508051906020019080838360005b8381101561213257818101518382015260200161211a565b50505050905090810190601f16801561215f5780820380516001836020036101000a031916815260200191505b5095505050505050602060405180830381600087803b15801561218157600080fd5b505af1158015612195573d6000803e3d6000fd5b505050506040513d60208110156121ab57600080fd5b506121b890508385612717565b6121c3836000612717565b5050505050505050565b6003546002146121dc57600080fd5b60018260ff16101580156121f4575060058260ff1611155b15156121ff57600080fd5b60018160ff1610158015612217575060058160ff1611155b151561222257600080fd5b6006546040517fa16e5dc2000000000000000000000000000000000000000000000000000000008152602060048201818152865160248401528651600160a060020a039094169363a16e5dc293889383926044909201919085019080838360005b8381101561229b578181015183820152602001612283565b50505050905090810190601f1680156122c85780820380516001836020036101000a031916815260200191505b5092505050602060405180830381600087803b1580156122e757600080fd5b505af11580156122fb573d6000803e3d6000fd5b505050506040513d602081101561231157600080fd5b5051151561231e57600080fd5b600654604080517fe3174a3e00000000000000000000000000000000000000000000000000000000815260048101918252865160448201528651600160a060020a039093169263e3174a3e928892889282916024810191606490910190602087019080838360005b8381101561239e578181015183820152602001612386565b50505050905090810190601f1680156123cb5780820380516001836020036101000a031916815260200191505b50838103825284518152845160209182019186019080838360005b838110156123fe5781810151838201526020016123e6565b50505050905090810190601f16801561242b5780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b15801561244c57600080fd5b505af1158015612460573d6000803e3d6000fd5b505050506040513d602081101561247657600080fd5b5051151561248357600080fd5b60048054604080517fe9fac80d0000000000000000000000000000000000000000000000000000000081523393810184815260248201928352885160448301528851600160a060020a039094169463e9fac80d9490938a939091606490910190602085019080838360005b838110156125065781810151838201526020016124ee565b50505050905090810190601f1680156125335780820380516001836020036101000a031916815260200191505b509350505050602060405180830381600087803b15801561255357600080fd5b505af1158015612567573d6000803e3d6000fd5b505050506040513d602081101561257d57600080fd5b5051151561258a57600080fd5b33600160a060020a0386161480156126f85750600654604080517f21f50b7700000000000000000000000000000000000000000000000000000000815260048101918252865160448201528651600160a060020a03909316926321f50b77928892889282916024810191606490910190602087019080838360005b8381101561261d578181015183820152602001612605565b50505050905090810190601f16801561264a5780820380516001836020036101000a031916815260200191505b50838103825284518152845160209182019186019080838360005b8381101561267d578181015183820152602001612665565b50505050905090810190601f1680156126aa5780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b1580156126cb57600080fd5b505af11580156126df573d6000803e3d6000fd5b505050506040513d60208110156126f557600080fd5b50515b156127105760ff828116908216141561271057600080fd5b5050505050565b6040805160018082528183019092526060918291829160208083019080388339505060408051600180825281830190925292955090506020808301908038833950506040805160018082528183019092529294509050602080830190803883390190505090503383600081518110151561278d57fe5b600160a060020a0390921660209283029091019091015281518590839060009081106127b557fe5b600160a060020a0390921660209283029091019091015280518490829060009081106127dd57fe5b90602001906020020181815250507fb2e1acee2fad7a352b9b70b4b0e6d2f18623c5110f671036a62b2819499b015f83838360405180806020018060200180602001848103845287818151815260200191508051906020019060200280838360005b8381101561285757818101518382015260200161283f565b50505050905001848103835286818151815260200191508051906020019060200280838360005b8381101561289657818101518382015260200161287e565b50505050905001848103825285818151815260200191508051906020019060200280838360005b838110156128d55781810151838201526020016128bd565b50505050905001965050505050505060405180910390a150505050505600a165627a7a7230582089712294a51fe6b20a47281ab72138c331a75bcee877e6e5c545aadb4884a9ba0029"

var testByteCode string = "608060405234801561001057600080fd5b50610346806100206000396000f3006080604052600436106100775763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630a9dfc51811461007c57806336720309146100a75780636d6b18af146100e8578063c9e77c88146100ff578063e35a50d81461011f578063fd1ee54c1461013f575b600080fd5b34801561008857600080fd5b506100936004610169565b604080519115158252519081900360200190f35b3480156100b357600080fd5b506100bf600435610180565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b3480156100f457600080fd5b5061009360046101b5565b34801561010b57600080fd5b5061009360048035602481019101356101c4565b34801561012b57600080fd5b5061009360048035602481019101356101db565b34801561014b57600080fd5b506101576004356101e9565b60408051918252519081900360200190f35b600061017781836002610208565b50600192915050565b600080548290811061018e57fe5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16905081565b60006101776001836002610285565b60006101d1818484610208565b5060019392505050565b60006101d160018484610285565b60018054829081106101f757fe5b600091825260209091200154905081565b828054828255906000526020600020908101928215610275579160200282015b8281111561027557815473ffffffffffffffffffffffffffffffffffffffff191673ffffffffffffffffffffffffffffffffffffffff843516178255602090920191600190910190610228565b506102819291506102cc565b5090565b8280548282559060005260206000209081019282156102c0579160200282015b828111156102c05782358255916020019190600101906102a5565b50610281929150610300565b6102fd91905b8082111561028157805473ffffffffffffffffffffffffffffffffffffffff191681556001016102d2565b90565b6102fd91905b8082111561028157600081556001016103065600a165627a7a723058200e8c68c33afa44615595055f617b887e38cbb67c6e69ac203bf7419620f212a80029"

var ownerAddr string = "0x0f768d36a18fff88ec6196dc7a82cb09cf4ebf47"
var ownerPriv string = "429bedad1bba8946d96b2d4365ec9c7d3c1085b7b7dd067f9d8631db3100a29d"

var sysAddr string = "0x10572de6d9eef252bb194f9c4b8055403a3357a2"
var sysPriv string = "a49a4f50b2dcf0af7237624281ebf6b5fc683db10a3be8941f1e885daac9814e"

var operAddr string = "0xc572ffd75174e24f95413ec1f0abd127948b36a5"
var operPriv string = "8404fcc511c885174fd59ecf7c297db2490ad989aeedf4c84a36c2190837d086"

var superAddr string = "0x65188459a1dc65984a0c7d4a397ed3986ed0c853"
var superPriv string = "7cb4880c2d4863f88134fd01a250ef6633cc5e01aeba4c862bedbf883a148ba8"

var testAddr1 string = "0xc409aaf73698fdb5995c4d85f6033d5e90d2f2bd"
var testPriv1 string = "64d18eb7061dff419581c1af98201b76c7ab6db538b1cf65123c470ccc6d5929"

const (
	accPriv  = "48deaa73f328f38d5fcb29d076b2b639c8491f97d245fc22e95a86366687903a"
	accAddr  = "28112ca022224ae7757bcd559666be5340ff109a"
	byteCode = `608060405234801561001057600080fd5b5061020c806100206000396000f30060806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806326f34b8b146100725780633c6bb436146100b35780633d4197f0146100de578063862c242e1461010b578063e1cb0e5214610142575b600080fd5b34801561007e57600080fd5b5061009d6004803603810190808035906020019092919050505061016d565b6040518082815260200191505060405180910390f35b3480156100bf57600080fd5b506100c861018d565b6040518082815260200191505060405180910390f35b3480156100ea57600080fd5b5061010960048036038101908080359060200190929190505050610193565b005b34801561011757600080fd5b50610140600480360381019080803590602001909291908035906020019092919050505061019d565b005b34801561014e57600080fd5b506101576101d7565b6040518082815260200191505060405180910390f35b600060016000838152602001908152602001600020600101549050919050565b60005481565b8060008190555050565b8160016000848152602001908152602001600020600001819055508060016000848152602001908152602001600020600101819055505050565b600080549050905600a165627a7a7230582090454576ac53e48f93db33e3502a6c1c9bff38697b8035a76500dc8ab84056b50029`
	//	abi            = `[{"constant": true,"inputs": [{"name": "_no","type": "uint256"}],"name": "getBatchVal","outputs": [{"name": "","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"},{"constant": true,"inputs": [],"name": "val","outputs": [{"name": "","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"},{"constant": false,"inputs": [{"name": "_val","type": "uint256"}],"name": "setVal","outputs": [],"payable": false,"stateMutability": "nonpayable","type": "function"},{"constant": false,"inputs": [{"name": "_no","type": "uint256"},{"name": "_val","type": "uint256"}],"name": "setBatchVal","outputs": [],"payable": false,"stateMutability": "nonpayable","type": "function"},{"constant": true,"inputs": [],"name": "getVal","outputs": [{"name": "","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"}]`
	caValidatorPri = "C4743CCC9CDDBE2E6F9F990D6255B06A779AAFE369336AA70A79672326045AFCC1F68967AB7F7FBCA19F64B7D329A68AB5B9BCA629BAE6931B41AD5C5A62BF66"
	addPeerPub     = "4B4375A2B7B2424A26E5F630CC458436F8391C941D46AA62EDAB614AE7A3A812"
	isCA           = true
	power          = 100
	blockHash      = "D824D437F02142FD92053C5039D182E4EE6AE704"
)

var testAddr2 string = "0x5f45ab0be3fc342e3a1033b293af0514d760ecf8"
var testPriv2 string = "d5236f9f29dfd15a8c6c42531f490b4682c6a39ef02f8de7ee2f0edab6037511"

var testAddr3 string = "0xc80bf8e2b390967dc908a4d37674473563036475"
var testPriv3 string = "fc70adae9f98412734748d9906e7524699b1e9d3eadf6b4f9b02df5fb911a4b9"

//var controlContractAddress = "0x6e649d44f531c64dd64a74559aec177e6e31ea86"
//var tokenContractAddress = "0xe7626be979385959384347fe0ad392859e69182f"
//var accountContractAddress = "0x496a99c449413cd1427a424df87ef26f6cad34c8"
//var transContractAddress = "0xd43891baf35eee671df6c461dba01a431de7917e"

var (
	newAddr, newPriv string
	nonce            uint64
	hash             string

	controlContractAddress string
	tokenContractAddress   string
	accountContractAddress string
	transContractAddress   string
	testContractAddress    string

	contractReceipt   string
	ccontractHash     string
	tcontractHash     string
	transcontractHash string
	acontractHash     string
	height            uint64

	client *sdk.GoSDK

	routines  int = 10
	sustained int = 100
)

func init() {
	client = sdk.New("127.0.0.1:46657", sdk.ZaCryptoType)
}

func GetNonce(source string) uint64 {
	nonce, _ := client.Nonce(source)
	return nonce
}

//func TestAddValidator(t *testing.T) {
//	err := client.AddValidator(caValidatorPri, addPeerPub, isCA, power)
//	if err != nil {
//		fmt.Println("add validator fail")
//	}
//}

//func TestRemoveValidator(t *testing.T) {
//	err := client.RemoveValidator(caValidatorPri, addPeerPub)
//	if err != nil {
//		fmt.Println("remove validator fail")
//	}
//}

//************************* 创建合约 *******************//
// control合约
func TestCreateControlContract(t *testing.T) {
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	var arg = sdk.ContractCreate{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:  controlAbis,
		Code: controlBytecode,
	}
	result, err := client.ContractCreate(&arg)
	client.ContractCreate(&arg)
	fmt.Println("======= create contract result:", result, err)
	assert.Nil(t, err)
	controlContractAddress = result["contract"].(string)
	ccontractHash = result["tx"].(string)
	t.Log(controlContractAddress, ccontractHash)
	fmt.Println("=====================查询Control合约地址：")
	t.Log(controlContractAddress)
	time.Sleep(time.Second * 3)
}

func TestQueryControlCreateReceipt(t *testing.T) {
	receipt, err := client.Receipt(ccontractHash)
	fmt.Println("=====================查询创建Control合约票据：")
	t.Log(receipt, err)
	assert.Nil(t, err)
	fmt.Println("==========:", receipt.From.String(), receipt.To.String(), receipt.Height, receipt.Timestamp, receipt.TxHash.String(), receipt.Status)
}

// token合约
func TestCreateTokenContract(t *testing.T) {
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	var arg = sdk.ContractCreate{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:  tokenAbis,
		Code: tokenBytecode,
	}
	result, err := client.ContractCreate(&arg)
	fmt.Println("======= create token contract result:", result, err)
	assert.Nil(t, err)
	tokenContractAddress = result["contract"].(string)
	tcontractHash = result["tx"].(string)
	t.Log(tokenContractAddress, tcontractHash)
	fmt.Println("=====================查询Token合约地址：")
	t.Log(tokenContractAddress)
	time.Sleep(time.Second * 3)
}

func TestQueryTokenCreateReceipt(t *testing.T) {
	receipt, err := client.Receipt(tcontractHash)
	fmt.Println("=====================查询创建Token合约票据：")
	t.Log(receipt, err)
	assert.Nil(t, err)
}

// account合约
func TestCreateAccountContract(t *testing.T) {
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	var arg = sdk.ContractCreate{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:  accountAbis,
		Code: accountBytecode,
	}
	result, err := client.ContractCreate(&arg)
	fmt.Println("======= create token contract result:", result, err)
	assert.Nil(t, err)
	accountContractAddress = result["contract"].(string)
	acontractHash = result["tx"].(string)
	t.Log(accountContractAddress, acontractHash)
	fmt.Println("=====================查询Account合约地址：")
	t.Log(accountContractAddress)
	time.Sleep(time.Second * 3)
}

func TestQueryAccountCreateReceipt(t *testing.T) {
	receipt, err := client.Receipt(acontractHash)
	fmt.Println("=====================查询创建Account合约票据：")
	t.Log(receipt, err)
	assert.Nil(t, err)
}

// trans合约
func TestCreateTransContract(t *testing.T) {
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	var arg = sdk.ContractCreate{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:  transAbis,
		Code: transBytecode,
	}
	result, err := client.ContractCreate(&arg)
	fmt.Println("======= create token contract result:", result, err)
	assert.Nil(t, err)
	transContractAddress = result["contract"].(string)
	transcontractHash = result["tx"].(string)
	t.Log(transContractAddress, transcontractHash)
	fmt.Println("=====================查询Trans合约地址：")
	t.Log(transContractAddress)
	time.Sleep(time.Second * 3)
}

func TestQueryTransCreateReceipt(t *testing.T) {
	receipt, err := client.Receipt(transcontractHash)
	fmt.Println("=====================查询创建Trans合约票据：")
	t.Log(receipt, err)
	assert.Nil(t, err)
}

//************************* 创建合约 *******************//

//************************* Control合约设置合约映射地址 *******************//
func TestExecuteControlContractSetAddress1(t *testing.T) {
	fmt.Println("=====================control合约配置account地址：")
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	params := []interface{}{"AccountContractAddress", accountContractAddress}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      controlAbis,
		Contract: controlContractAddress,
		Method:   "changeContractAddress",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
	fmt.Println("==========:", receipt.From.String(), receipt.To.String(), receipt.Height, receipt.Timestamp, receipt.TxHash.String(), receipt.Status)
}

func TestExecuteControlContractSetAddress2(t *testing.T) {
	fmt.Println("=====================control合约配置assert地址：")
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	params := []interface{}{"AssetTokenAddress", tokenContractAddress}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      controlAbis,
		Contract: controlContractAddress,
		Method:   "changeContractAddress",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
	fmt.Println("==========:", receipt.From.String(), receipt.To.String(), receipt.Height, receipt.Timestamp, receipt.TxHash, receipt.Status)
}

func TestQueryCotrolContractAccountAddress(t *testing.T) {
	fmt.Println("=====================查询设置的Account合约地址:")
	params := []interface{}{"AccountContractAddress"}
	ownerNonce := GetNonce(ownerAddr)
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      controlAbis,
		Contract: controlContractAddress,
		Method:   "queryContractAddress",
		Params:   params,
	}

	resp, err := client.ContractRead(&arg)
	fmt.Println("======= read contract:", resp.(common.Address).String())
	assert.Nil(t, err)
}

func TestQueryCotrolContractTokenAddress(t *testing.T) { // 0
	fmt.Println("=====================查询设置的Token合约地址:")
	params := []interface{}{"AssetTokenAddress"}
	ownerNonce := GetNonce(ownerAddr)
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      controlAbis,
		Contract: controlContractAddress,
		Method:   "queryContractAddress",
		Params:   params,
	}

	resp, err := client.ContractRead(&arg)
	fmt.Println("======= read contract:", resp.(common.Address).String())
	assert.Nil(t, err)
}

//************************* Control合约设置合约映射地址 *******************//

//************************* 设置管理账户 *******************//
func TestExecuteTransContractSetSysAccount(t *testing.T) {
	fmt.Println("=====================trans合约设置sys管理账户：")
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	params := []interface{}{ownerAddr}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      transAbis,
		Contract: transContractAddress,
		Method:   "setSysAccount",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
}

func TestExecuteTransContractSetOperaAccount(t *testing.T) {
	fmt.Println("=====================trans合约设置operation管理账户：")
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	params := []interface{}{ownerAddr}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      transAbis,
		Contract: transContractAddress,
		Method:   "setOperationAccount",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
}

func TestQueryTransContractAccounts(t *testing.T) { // 0
	time.Sleep(time.Second * 3)
	fmt.Println("=====================查询trans合约管理账户列表：")
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   0,
		},
		ABI:      transAbis,
		Contract: transContractAddress,
		Method:   "queryCurrentAccounts",
		Params:   nil,
	}

	resp, err := client.ContractRead(&arg)
	fmt.Println("======= read contract:", resp)
	assert.Nil(t, err)
}

func TestExecuteToenContractSetSysAccount(t *testing.T) {
	fmt.Println("=====================token合约设置sys管理账户：")
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	params := []interface{}{ownerAddr}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      tokenAbis,
		Contract: tokenContractAddress,
		Method:   "setSysAccount",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
}

func TestExecuteTokenContractSetOperaAccount(t *testing.T) {
	fmt.Println("=====================token合约设置operation管理账户：")
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	params := []interface{}{ownerAddr}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      tokenAbis,
		Contract: tokenContractAddress,
		Method:   "setOperationAccount",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
}

func TestQueryTokenContractAdminAccounts(t *testing.T) { // 0
	time.Sleep(time.Second * 3)
	fmt.Println("=====================查询token合约管理账户列表：")
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   0,
		},
		ABI:      tokenAbis,
		Contract: tokenContractAddress,
		Method:   "queryCurrentAccounts",
		Params:   nil,
	}

	resp, err := client.ContractRead(&arg)
	fmt.Println("======= read contract:", resp)
	assert.Nil(t, err)
}

func TestExecuteAcountContractSetSysAccount(t *testing.T) {
	fmt.Println("=====================Account合约设置sys管理账户：")
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	params := []interface{}{ownerAddr}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      accountAbis,
		Contract: accountContractAddress,
		Method:   "setSysAccount",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
}

func TestExecuteAcountContractSetOpeartionAccount(t *testing.T) {
	fmt.Println("=====================Account合约设置operation管理账户：")
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	params := []interface{}{ownerAddr}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      accountAbis,
		Contract: accountContractAddress,
		Method:   "setOperationAccount",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
}

func TestQueryAccountContractAdminAccounts(t *testing.T) { // 0
	time.Sleep(time.Second * 3)
	fmt.Println("=====================查询Account合约管理账户列表：")
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   0,
		},
		ABI:      accountAbis,
		Contract: accountContractAddress,
		Method:   "queryCurrentAccounts",
		Params:   nil,
	}

	resp, err := client.ContractRead(&arg)
	fmt.Println("======= read contract:", resp)
	assert.Nil(t, err)
}

//************************* 设置管理账户 *******************//

//************************* 初始化合约 *******************//
func TestExecuteAccountContractInit(t *testing.T) {
	fmt.Println("=====================account合约初始化")
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	params := []interface{}{controlContractAddress}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      accountAbis,
		Contract: accountContractAddress,
		Method:   "initContracts",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
}

func TestExecuteTrnsContractInit(t *testing.T) {
	fmt.Println("=====================trans合约初始化")
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	params := []interface{}{controlContractAddress}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      transAbis,
		Contract: transContractAddress,
		Method:   "initContracts",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
}

//************************* 初始化合约 *******************//

//************************* 设置Accessible权限 *******************//
func TestExecuteTokenContractSetAccessAccount1(t *testing.T) {
	fmt.Println("=====================token合约设置account合约Accessible权限")
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	params := []interface{}{accountContractAddress, true}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      tokenAbis,
		Contract: tokenContractAddress,
		Method:   "setAccessibleAccounts",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
}

func TestExecuteTokenContractSetAccessAccount2(t *testing.T) {
	fmt.Println("=====================token合约设置trans合约Accessible权限")
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	params := []interface{}{transContractAddress, true}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      tokenAbis,
		Contract: tokenContractAddress,
		Method:   "setAccessibleAccounts",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
}

func TestExecuteTokenContractSetAccessAccount3(t *testing.T) {
	fmt.Println("=====================token合约设置owner拥有Accessible权限")
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	params := []interface{}{ownerAddr, true}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      tokenAbis,
		Contract: tokenContractAddress,
		Method:   "setAccessibleAccounts",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
}

func TestQueryTokenContractAccessAccounts(t *testing.T) { // 0
	fmt.Println("查询token合约access accounts是否有权限：")
	params1 := []interface{}{ownerAddr}
	var arg1 = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   0,
		},
		ABI:      tokenAbis,
		Contract: tokenContractAddress,
		Method:   "_accessibleAccounts",
		Params:   params1,
	}

	resp1, err := client.ContractRead(&arg1)
	fmt.Println("======= read contract:", resp1)
	assert.Nil(t, err)

	params2 := []interface{}{transContractAddress}
	var arg2 = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   0,
		},
		ABI:      tokenAbis,
		Contract: tokenContractAddress,
		Method:   "_accessibleAccounts",
		Params:   params2,
	}

	resp2, err := client.ContractRead(&arg2)
	fmt.Println("======= read contract:", resp2)
	assert.Nil(t, err)

	params3 := []interface{}{accountContractAddress}
	var arg3 = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   0,
		},
		ABI:      tokenAbis,
		Contract: tokenContractAddress,
		Method:   "_accessibleAccounts",
		Params:   params3,
	}

	resp3, err := client.ContractRead(&arg3)
	fmt.Println("======= read contract:", resp3)
	assert.Nil(t, err)
}

func TestAcountExecuteContractConfigureAccessAccounts1(t *testing.T) {
	fmt.Println("=====================account合约设置owner拥有Accessible权限")
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	params := []interface{}{ownerAddr, true}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      accountAbis,
		Contract: accountContractAddress,
		Method:   "configureAccessibleAccounts",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
}

func TestAcountExecuteContractConfigureAccessAccounts2(t *testing.T) {
	fmt.Println("=====================account合约设置trans合约Accessible权限")
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	params := []interface{}{transContractAddress, true}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      accountAbis,
		Contract: accountContractAddress,
		Method:   "configureAccessibleAccounts",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
}

func TestAcountExecuteContractConfigureAccessAccounts3(t *testing.T) {
	fmt.Println("=====================account合约设置token合约Accessible权限")
	ownerNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownerNonce)
	params := []interface{}{tokenContractAddress, true}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownerNonce,
		},
		ABI:      accountAbis,
		Contract: accountContractAddress,
		Method:   "configureAccessibleAccounts",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
}

func TestQueryAccountContractAccessAccounts(t *testing.T) { // 0
	fmt.Println("=====================查询account合约账户access权限")
	params1 := []interface{}{ownerAddr}
	var arg1 = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   0,
		},
		ABI:      accountAbis,
		Contract: accountContractAddress,
		Method:   "_accessibleAccounts",
		Params:   params1,
	}

	resp1, err := client.ContractRead(&arg1)
	fmt.Println("======= read contract:", resp1)
	assert.Nil(t, err)

	params2 := []interface{}{tokenContractAddress}
	var arg2 = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   0,
		},
		ABI:      accountAbis,
		Contract: accountContractAddress,
		Method:   "_accessibleAccounts",
		Params:   params2,
	}

	resp2, err := client.ContractRead(&arg2)
	fmt.Println("======= read contract:", resp2)
	assert.Nil(t, err)

	params3 := []interface{}{transContractAddress}
	var arg3 = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   0,
		},
		ABI:      accountAbis,
		Contract: accountContractAddress,
		Method:   "_accessibleAccounts",
		Params:   params3,
	}

	resp3, err := client.ContractRead(&arg3)
	fmt.Println("======= read contract:", resp3)
	assert.Nil(t, err)
}

//************************* 设置Accessible权限 *******************//

//************************* transfer发行资产 *******************//
func TestExecuteTrnsContractPublish(t *testing.T) {
	fmt.Println("=====================trans合约给testAddr1发行资产")
	ownNonce := GetNonce(ownerAddr)
	fmt.Println("nonce:", ownNonce)
	params := []interface{}{201904111446, testAddr1, "ZABT", "T20190411", 1000000, 500000000000000, 400000000000000, 1}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   ownNonce,
		},
		ABI:      transAbis,
		Contract: transContractAddress,
		Method:   "publish",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)

	abiJson, err := abi.JSON(strings.NewReader(transAbis))
	if err != nil {
		t.Log(err)
	}
	l := sdk.EventLog{(types.Receipt)(*receipt), abiJson}
	actual := l.Print("event_publish")
	t.Log(actual)
}

//************************* transfer发行资产 *******************//

//************************* transfer合约转账 *******************//
func TestTransExecuteContractTransfer(t *testing.T) {
	fmt.Println("=====================trans合约testAddr1给testAddr2转账标准资产1")
	testNonce := GetNonce(testAddr1)
	fmt.Println("nonce:", testNonce)
	params := []interface{}{1, "0x383c9ed0cb41bd9f72eb5e9975071913b4fce7c0", "T20190411", "T20190411", 2, 2, 300000000000000}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: testPriv1,
			Nonce:   testNonce,
		},
		ABI:      transAbis,
		Contract: transContractAddress,
		Method:   "transfer",
		Params:   params,
	}
	result, err := client.ContractCall(&arg)
	assert.Nil(t, err)
	t.Log(result)

	time.Sleep(time.Second * 3)
	receipt, err := client.Receipt(result)
	t.Log(receipt)
}

//************************* transfer合约转账 *******************//

//************************* 查询资产信息 *******************//
func TestQueryTokenContractTokenInfo(t *testing.T) { // [ZABT T20190411 [196 9 170 247 54 152 253 181 153 92 77 133 246 3 61 94 144 210 242 189] [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] 0 0 0 0  1564049034 1000000]
	fmt.Println("=====================查询token合约资产info:")
	params := []interface{}{"T20190411"}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   0,
		},
		ABI:      tokenAbis,
		Contract: tokenContractAddress,
		Method:   "queryTokenInfo",
		Params:   params,
	}

	resp, err := client.ContractRead(&arg)
	fmt.Println("======= read contract:", resp)
	assert.Nil(t, err)
}

func TestQueryToenContractChkSymbolToken(t *testing.T) { // true
	fmt.Println("=====================查询token合约资产symbol info:")
	params := []interface{}{"ZABT", "T20190411"}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   0,
		},
		ABI:      tokenAbis,
		Contract: tokenContractAddress,
		Method:   "chkSymbolToken",
		Params:   params,
	}

	resp, err := client.ContractRead(&arg)
	fmt.Println("======= read contract:", resp)
	assert.Nil(t, err)
}

func TestQueryAccountContractTokenInfo(t *testing.T) { // [100000000000000 100000000000000 0 0 0]
	fmt.Println("=====================查询account testAddr1 资产info:")
	params := []interface{}{testAddr1, "ZABT", "T20190411"}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: ownerPriv,
			Nonce:   0,
		},
		ABI:      accountAbis,
		Contract: accountContractAddress,
		Method:   "queryAccountTokenInfo",
		Params:   params,
	}

	resp, err := client.ContractRead(&arg)
	fmt.Println("======= read contract:", resp)
	assert.Nil(t, err)
}

func TestQueryAccountContractTokenTypeAmount1(t *testing.T) { // 300000000000000
	fmt.Println("=====================查询account testAddr1 标准资产info:")
	params := []interface{}{"0x383c9ed0cb41bd9f72eb5e9975071913b4fce7c0", "T20190411", 2}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: testPriv2,
			Nonce:   0,
		},
		ABI:      accountAbis,
		Contract: accountContractAddress,
		Method:   "queryAccountTokenTypeAmount",
		Params:   params,
	}

	resp, err := client.ContractRead(&arg)
	fmt.Println("======= read contract:", resp)
	assert.Nil(t, err)
}

func TestQueryAccountContractTokenTypeAmount2(t *testing.T) { // 100000000000000
	fmt.Println("=====================查询account testAddr2 标准资产info:")
	params := []interface{}{testAddr1, "T20190411", 2}
	var arg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: testPriv1,
			Nonce:   0,
		},
		ABI:      accountAbis,
		Contract: accountContractAddress,
		Method:   "queryAccountTokenTypeAmount",
		Params:   params,
	}

	resp, err := client.ContractRead(&arg)
	fmt.Println("======= read contract:", resp)
	assert.Nil(t, err)
}
