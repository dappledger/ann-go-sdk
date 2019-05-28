package ikhofi

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dappledger/AnnChain-go-sdk/db"
	"github.com/dappledger/AnnChain-go-sdk/wire"
)

const (
	databaseCache   = 128
	databaseHandles = 1024
)

var (
	lastBlockKey        = []byte("lastblock")
	ErrRevertFromBackup = errors.New("revert from backup,not find data")
	ErrDataTransfer     = errors.New("data transfer err")
	ErrBranchNameUsed   = errors.New("app:branch name has been used")
)

type BaseApplication struct {
	Database         db.DB
	InitializedState bool
}

// InitBaseApplication must be the first thing to be called when an application embeds BaseApplication
func (ba *BaseApplication) InitBaseApplication(name string, datadir string) (err error) {
	if ba.Database, err = db.NewGoLevelDB(name, datadir); err != nil {
		return err
	}
	ba.InitializedState = true
	return nil
}

func (ba *BaseApplication) LoadLastBlock(t interface{}) (res interface{}, err error) {
	buf := ba.Database.Get(lastBlockKey)
	if len(buf) != 0 {
		r, n, err := bytes.NewReader(buf), new(int), new(error)
		res = wire.ReadBinaryPtr(t, r, 0, n, err)
		if *err != nil {
			return nil, *err
		}
	} else {
		return nil, errors.New("empty")
	}
	return
}

func (ba *BaseApplication) SaveLastBlockByKey(key []byte, lastBlock interface{}) {
	buf, n, err := new(bytes.Buffer), new(int), new(error)
	wire.WriteBinary(lastBlock, buf, n, err)
	if *err != nil {
		panic(*err)
	}
	ba.Database.SetSync(key, buf.Bytes())
}

type AppTool struct {
	BaseApplication
	lastBlock LastBlockInfo
}

func (t *AppTool) Init(datadir string) error {
	if err := t.InitBaseApplication(APP_NAME, datadir); err != nil {
		return err
	}
	ret, err := t.LoadLastBlock(&t.lastBlock)
	if err != nil {
		return err
	}
	tmp, ok := ret.(*LastBlockInfo)
	if !ok {
		return ErrDataTransfer
	}
	t.lastBlock = *tmp
	return nil
}

func (t *AppTool) LastHeightHash() (int64, []byte) {
	return int64(t.lastBlock.Height), t.lastBlock.AppHash
}

func (t *AppTool) BackupLastBlock(branchName string) error {
	return t.backupLastBlockData(branchName, &t.lastBlock)
}

func (t *AppTool) backupLastBlockData(branchName string, lastBlock interface{}) error {
	preKeyName := []byte(fmt.Sprintf("%s-%s", lastBlockKey, branchName))
	dataBs := t.Database.Get(preKeyName)
	if len(dataBs) > 0 {
		return ErrBranchNameUsed
	}
	t.SaveLastBlockByKey(preKeyName, lastBlock)
	return nil
}

func (t *AppTool) SaveNewLastBlock(fromHeight int64, fromAppHash []byte) error {
	newBranchBlock := LastBlockInfo{
		Height:  fromHeight,
		AppHash: fromAppHash,
	}
	t.SaveLastBlockByKey(lastBlockKey, newBranchBlock)
	// TODO
	return nil
}
