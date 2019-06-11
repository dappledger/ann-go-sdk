package sdk

type CyrptoType string

const (
	ZaCryptoType CyrptoType = "ZA"
)

type CommitType string

const (
	TypeSyn  CommitType = "syn"
	TypeAsyn CommitType = "asyn"
)

type GoSDK struct {
	rpcAddr    string
	cryptoType CyrptoType
}

func (gs *GoSDK) Url() string {
	return gs.rpcAddr
}

func New(rpcAddr string, cryptoType CyrptoType) *GoSDK {
	return &GoSDK{
		rpcAddr,
		cryptoType,
	}
}

func (gs *GoSDK) JsonRPCCall(method string, params []byte, result interface{}) error {
	return gs.sendTxCall(method, params, result)
}

func (gs *GoSDK) Put(privKey string, value []byte, typ CommitType) (string, error) {
	return gs.put(privKey, value, typ)
}

func (gs *GoSDK) Get(key string) ([]byte, error) {
	return gs.get(key)
}

func (gs *GoSDK) AccountCreate() (Account, error) {
	return gs.accountCreate()
}
