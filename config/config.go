package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var ConfigFile = ".godo/config.json"

type Config struct {
	DataFile       string `json:"DataFile"`
	PrivateKeyFile string `json:"PrivateKeyFile"`
	PublicKeyFile  string `json:"PublicKeyFile"`
	GodoServerUrl  string `json:"GodoServerUrl"`
}

func init() {
	initDataFile()
}

var ClientConfig = &Config{}

func initDataFile() {
	homeDir := os.Getenv("HOME")
	configFile := homeDir + "/" + ConfigFile
	if !pathExist(configFile) {
		fmt.Printf("config file:$HOME/%s do not exist\n", ConfigFile)
		os.Exit(CONFIG_FILE_DO_NOT_EXIST)
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

	if err := json.Unmarshal(configData, ClientConfig); err != nil {
		panic(err)
	}
}

func pathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
