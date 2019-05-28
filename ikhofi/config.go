package ikhofi

import (
	"os"
	"path"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"fmt"
	"io/ioutil"
	"path/filepath"
)

const (
	RUNTIME_ENV     = "ANN_RUNTIME"
	IKHOFI_PATH     = "IKHOFI_PATH"
	DEFAULT_RUNTIME = ".ann_runtime"
	DATADIR         = "contract_data"
	DBDATADIR       = "chaindata"
	CONFIGFILE      = "ikhofi.yaml"
)

type IkhofiConfig struct {
	Db struct {
		DbType            string `yaml:"type"`
		DbPath            string `yaml:"path"`
		CacheSize         int    `yaml:"cacheSize"`
		DestroyAfterClose bool   `yaml:"destroyAfterClose"`
	}
}

func runtimeDir(root string) string {
	if root != "" {
		if root[0:1] != "/" {
			pwd, _ := os.Getwd()
			root = filepath.Join(pwd, root)
		}
		return root
	}
	if runtimePath, exists := os.LookupEnv(RUNTIME_ENV); exists {
		return runtimePath
	}
	return path.Join(os.Getenv("HOME"), DEFAULT_RUNTIME)
}

func getYamlBytes(dbPath string) (yamlBytes []byte, err error) {
	cfg := IkhofiConfig{}
	err = yaml.Unmarshal([]byte(CONFIGTPL), &cfg)
	if err != nil {
		return
	}
	cfg.Db.DbType = "leveldb"
	cfg.Db.DbPath = dbPath
	cfg.Db.CacheSize = 67108864
	cfg.Db.DestroyAfterClose = false
	yamlBytes, err = yaml.Marshal(&cfg)
	return
}

func EnsureDir(dir string, mode os.FileMode) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, mode)
		if err != nil {
			return fmt.Errorf("Could not create directory %v. %v", dir, err)
		}
	}
	return nil
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func WriteFile(filePath string, contents []byte, mode os.FileMode) error {
	err := ioutil.WriteFile(filePath, contents, mode)
	if err != nil {
		return err
	}
	// fmt.Printf("File written to %v.\n", filePath)
	return nil
}

func MustWriteFile(filePath string, contents []byte, mode os.FileMode) {
	err := WriteFile(filePath, contents, mode)
	if err != nil {
		Exit(fmt.Sprintf("MustWriteFile failed: %v", err))
	}
}

func Exit(s string) {
	fmt.Printf(s + "\n")
	os.Exit(1)
}

func initRuntime(root string) string {
	EnsureDir(root, 0700)
	EnsureDir(path.Join(root, DATADIR), 0700)
	configFilePath := path.Join(root, CONFIGFILE)
	if !FileExists(configFilePath) {
		yamlBytes, err := getYamlBytes(path.Join(root, DATADIR, DBDATADIR, "/"))
		if err != nil {
			Exit("can not generate ikhofi.yaml file.")
		}
		MustWriteFile(configFilePath, yamlBytes, 0644)
	}
	return configFilePath
}

func initIkhofiConfig(conf *viper.Viper) *viper.Viper {
	runtime := runtimeDir(conf.GetString("db_dir"))
	configFilePath := initRuntime(runtime)

	conf.SetDefault("ikhofi_config", configFilePath)

	return conf
}
