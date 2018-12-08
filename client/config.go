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
	DataFile      string       `json:"DataFile"`
	RsaConfig     RsaConfig    `json:"RsaConfig"`
	AesGCMConfig  AesGCMConfig `json:"AesGCMConfig"`
	GodoServerUrl string       `json:"GodoServerUrl"`
	GithubRepo    string       `json:"GithubRepo"`
}

type RsaConfig struct {
	PublicKeyFile  string `json:"PublicKeyFile"`
	PrivateKeyFile string `json:"PrivateKeyFile"`
}

type AesGCMConfig struct {
	Key   string `json:"Key"`
	Nonce string `json:"Nonce"`
}

var ClientConfig = initDataFile()

func initDataFile() *Config {
	homeDir := os.Getenv("HOME")
	configFile := homeDir + "/" + ConfigFile
	if !pathExist(configFile) {
		fmt.Printf("consts file:$HOME/%s do not exist\n", ConfigFile)
		os.Exit(consts.ConfigFileDoNotExist)
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
