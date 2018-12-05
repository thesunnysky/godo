package godo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/thesunnysky/godo/consts"
	"io/ioutil"
	"os"
)

var ConfigFile = ".godo/config.json"

type Config struct {
	DataFile      string `json:"DataFile"`
	AesGCMKey     string `json:"AesGCMKey"`
	AesGCMNonce   string `json:"AesGCMNonce"`
	GodoServerUrl string `json:"GodoServerUrl"`
}

var ClientConfig = initDataFile()

func initDataFile() *Config {
	homeDir := os.Getenv("HOME")
	configFile := homeDir + "/" + ConfigFile
	if !pathExist(configFile) {
		fmt.Printf("consts file:$HOME/%s do not exist\n", ConfigFile)
		os.Exit(consts.CONFIG_FILE_DO_NOT_EXIST)
	}
	f, err := os.Open(configFile)
	if err != nil {
		panic(nil)
	}
	defer f.Close()

	configData, err := ioutil.ReadAll(bufio.NewReader(f))
	if err != nil {
		panic(err)
	}

	ClientConfig := &Config{}
	if err := json.Unmarshal(configData, ClientConfig); err != nil {
		panic(err)
	}
	return ClientConfig
}

func pathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
