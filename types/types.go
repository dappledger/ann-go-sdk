package types

type QueryType = byte

const (
	QueryType_Contract  QueryType = 0
	QueryType_Nonce     QueryType = 1
	QueryType_Balance   QueryType = 2
	QueryType_Receipt   QueryType = 3
	QueryType_Existence QueryType = 4
	QueryType_PayLoad   QueryType = 5
)
