package ikhofi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/dappledger/AnnChain-go-sdk/common"
	agtypes "github.com/dappledger/AnnChain-go-sdk/types"
)

type IKHOFIApp struct {
	BaseApplication

	core interface{}

	datadir  string
	stateMtx sync.Mutex

	Config *viper.Viper
}

type LastBlockInfo struct {
	Height  int64
	AppHash []byte
}

type StartParams struct {
	ConfigPath   string
	StateHashHex string
}

var (
	ikhofiSigner = DawnSigner{}

	SystemContractId           = "system"
	SystemDeployMethod         = "deploy"
	SystemUpgradeMethod        = "upgrade"
	SystemQueryContractIdExits = "contract"

	APP_NAME = "ikhofi"

	serverUrl = ""
)

func NewIKHOFIApp(config *viper.Viper) (*IKHOFIApp, error) {
	config = initIkhofiConfig(config)

	serverUrl = config.GetString("ikhofi_laddr")
	if serverUrl == "" {
		return nil, errors.Wrap(errors.Errorf("app miss configuration ikhofi_laddr"), "app error")
	}
	// TODO
	if serverUrl[0:6] != "http://" {
		serverUrl = "http://" + serverUrl
	}

	app := IKHOFIApp{
		datadir: config.GetString("db_dir"),

		Config: config,
	}

	var err error
	if err = app.BaseApplication.InitBaseApplication(APP_NAME, app.datadir); err != nil {
		return nil, errors.Wrap(err, "app error")
	}

	return &app, nil
}

func (app *IKHOFIApp) Start() (err error) {
	lastBlock := &LastBlockInfo{
		Height:  0,
		AppHash: make([]byte, 0),
	}

	trieRoot := ""
	if res, err := app.LoadLastBlock(lastBlock); err == nil && res != nil {
		lastBlock = res.(*LastBlockInfo)
	}
	if err != nil {
		//log.Error("fail to load last block", zap.Error(err))
		return
	}
	trieRoot = common.Bytes2Hex(lastBlock.AppHash)

	path := app.Config.GetString("ikhofi_config")
	if path[0:1] != "/" {
		pwd, _ := os.Getwd()
		path = filepath.Join(pwd, path)
	}

	startParams := StartParams{
		path,
		trieRoot,
	}

	startParamsB, err := json.Marshal(startParams)
	if err != nil {
		return
	}
	resultB, err := app.post("start", "application/json", startParamsB)
	result, _ := strconv.ParseBool(string(resultB))
	if !result {
		app.Stop()
		return errors.Wrap(errors.Errorf("ikhofi jvm start error"), "[IKHOFIApp Start]")
	}

	return nil
}

func (app *IKHOFIApp) Stop() {
	app.get("stop")
	return
}

func (app *IKHOFIApp) CompatibleWithAngine() {}

func (app *IKHOFIApp) Query(query []byte) agtypes.Result {
	txpb := &TransactionPb{}
	if err := proto.Unmarshal(query, txpb); err != nil {
		return agtypes.NewError(agtypes.CodeType_BaseInvalidInput, err.Error())
	}
	tx := Pb2Transaction(txpb)

	queryData := &Query{
		Version: 1,
		Id:      tx.To,
		Method:  tx.Method,
	}
	if len(tx.Args) > 0 {
		queryData.Args = tx.Args
	}

	queryDataBytes, _ := proto.Marshal(queryData)
	res, err := app.post("query", "text/plain", queryDataBytes)
	if err != nil {
		return agtypes.NewError(agtypes.CodeType_BaseInvalidOutput, err.Error())
	}

	return agtypes.NewResultOK(res, "")
}

func (app *IKHOFIApp) SetCore(core interface{}) {
	app.core = core
}

func (app *IKHOFIApp) get(method string) (result []byte, err error) {
	url := serverUrl + "/" + method
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}

func (app *IKHOFIApp) post(method, contentType string, data []byte) (result []byte, err error) {
	url := serverUrl + "/" + method
	resp, err := http.Post(url, contentType, bytes.NewBuffer(data))
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return
	}

	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}
