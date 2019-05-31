# 依赖库：

## github.com

- github.com/fsnotify/fsnotify
- github.com/golang/protobuf/proto
- github.com/golang/snappy
- github.com/hashicorp/hcl
- github.com/magiconair/properties
- github.com/mitchellh/mapstructure
- github.com/pelletier/go-toml
- github.com/pkg/errors
- github.com/spf13/viper
- github.com/syndtr/goleveldb/leveldb

## golang.org

- golang.org/x/crypto
- golang.org/x/text

## gopkg.in

- gopkg.in/yaml.v2

# 使用说明

## 1、初始化

构造sdk对象;其他相关操作,都是通过sdk对象实现;

```golang
func NewGoSDK(rpcaddr, crypto_type string) *GoSDK
```

参数说明：
- rpcaddr：区块链rpc接口地址;
- crypto_type：加密类型;
  - ZA:默认加密类型,包括ed25519/secp256k/Keccak256;
  - GM:国密加密类型,开源版暂不支持,仅商业版支持;
  

## 链健康检查

```golang
func (gs *GoSDK) CheckHealth() (bool, error) 
```

输出参数：
- bool:是否成功;
- error：发生错误;


## 生成账号

```golang
func (gs *GoSDK) AccountGenerate() (*Account, error)
```

输出参数：

- Account:账号结构体;
- error：发生错误;

```golang
type Account struct {
	Privkey []byte `json:"privkey"`
	Address []byte `json:"address"`
}
```

## 查询nonce

```golang
func (gs *GoSDK) Nonce(addr string) (uint64, error)
```

输入参数：
- addr:查询地址;

输出参数：
- uint64:查询到的nonce（error为nil有效）;
- error：发生错误;



## 创建合约：

```golang
func (gs *GoSDK) ContractCreate(contract *ContractCreate) (map[string]interface{}, error)
```

输入参数：

```golang
type AccountBase struct {
	PrivKey string `json:"privkey"`//私钥
	Nonce   uint64 `json:"nonce"`  
}
type ContractCreate struct {
	AccountBase
	Code   string        `json:"code"`          //合约编译之后的字节码;
	ABI    string        `json:"abiDefinition"` //合约的ABI;
	Params []interface{} `json:"params"`        //合约构造函数用到的参数
}
```

输出参数：

- error:如果调用出错,会返回对应的错误;
- map[string]interface{}： 本交易hash和合约地址;例子如下：

```golang
ret := map[string]interface{}{
		"tx":       hex.EncodeToString(hash),
		"contract": contractAddr.Hex(),
	}
```

## 调用合约(写)：

```golang
//同步调用
func (gs *GoSDK) ContractCall(contractMethod *ContractMethod) (string, error)
```
```golang
//异步调用
func (gs *GoSDK) ContractAsync(contractMethod *ContractMethod) (string, error)
```

输入参数：

```golang
type ContractMethod struct {
	AccountBase
	Contract string        `json:"contract"`
	ABI      string        `json:"abiDefinition"`   //合约的ABI;
	Method   string        `json:"method"`          //要调用的合约函数 的名称
	Params   []interface{} `json:"params"`          //要调用的合约函数 的参数
}
```

输出参数：

- string: 本操作hash和合约地址;
- error:  如果调用出错,会返回对应的错误;



## 调用合约(读)：

```golang
func (gs *GoSDK) ContractRead(contractMethod *ContractMethod) (interface{}, error) 
```

输入参数：

```golang
type ContractMethod struct {
	AccountBase
	Contract string        `json:"contract"`
	ABI      string        `json:"abiDefinition"`   //合约的ABI;
	Method   string        `json:"method"`          //要调用的合约函数 的名称
	Params   []interface{} `json:"params"`          //要调用的合约函数 的参数
}
```

输出参数：

- interface{}:  本操作具体读取到的数据;
- error:		如果调用出错,会返回对应的错误;




## 获取收据 

```golang
func (gs *GoSDK) Receipt(hash string) (*types.ReceiptForStorage, error) 
```

输入参数：
- hash: 交易hash;

输出参数：
- ReceiptForStorage:收据内容;
- error：发生错误;

## 根据块hash获取交易列表 

```golang
func (gs *GoSDK) Block(hash string) ([]string, int, error)
```

输入参数：
- hash: 区块hash;

输出参数：
- []string: 交易hash列表;
- int: 交易数量;
- error：发生错误;

## 已签名交易

上述交易内置sdk进行签名了,所以需要带上私钥（用于签名）。  
这里传的参数,是直接已经转换为已经签名过的字节流,直接发送给链就可以了。  


```golang
//同步调用
func (gs *GoSDK) TranscationSignature(tx string) (string, error)
```
//异步调用
```golang
func (gs *GoSDK) TranscationSignatureAsync(tx string) (string, error)
```

输入参数：

- tx：已经签名过的序列化过的请求参数;

输出参数：
- string: 本交易hash;
- error:  如果调用出错,会返回对应的错误;


