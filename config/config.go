package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const CONFIG_FILE = ".godo/config.json"

const INVALID_PARAMETER_VALUE = 1

const CONFIG_FILE_DO_NOT_EXIST = 2

const FILE_MAKS = 0666

const LINE_SEPARATOR = '\n'

const DEFAULT_LINE_CACHE = 50

type Config struct {
	DataFile       string `json:"DataFile"`
	PrivateKeyFile string `json:"PrivateKeyFile"`
	PublicKeyFile  string `json:"PublicKeyFile"`
}

func init() {
	initDataFile()
}

var CONF = &Config{}

func initDataFile() {
	homeDir := os.Getenv("HOME")
	configFile := homeDir + "/" + CONFIG_FILE
	if !pathExist(configFile) {
		fmt.Printf("config myfile:$HOME/%s do not exist\n", CONFIG_FILE)
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

	if err := json.Unmarshal(configData, CONF); err != nil {
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
