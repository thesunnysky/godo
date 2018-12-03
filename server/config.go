package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

//godo-server consts
const (
	GODO_DATA_FILE = "GodoDataFile"
)

const CONFIG_FILE = ""

type Config struct {
}

func NewConfig() (*Config, error) {
	config := &Config{}
	homeDir := os.Getenv("HOME")
	configFile := homeDir + "/" + CONFIG_FILE
	if !pathExist(configFile) {
		str := fmt.Sprintf("config myfile:$HOME/%s do not exist\n", CONFIG_FILE)
		return nil, errors.New(str)
	}
	f, err := os.Open(configFile)
	if err != nil {
	}
	defer f.Close()

	configData, err := ioutil.ReadAll(bufio.NewReader(f))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(configData, config); err != nil {
		return nil, err
	}
	return config, nil
}

func pathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
