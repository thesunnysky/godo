package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thesunnysky/godo/util"
	"io/ioutil"
	"os"
)

const ConfigFile = ".godo/server.json"

type Config struct {
	DataDir string `json:"DataDir"`
}

var ServerConfig = &Config{}

func init() {
	_ = initConfig()
}

func initConfig() error {
	homeDir := os.Getenv("HOME")
	configFile := homeDir + "/" + ConfigFile
	if !util.PathExist(configFile) {
		str := fmt.Sprintf("config file:$HOME/%s do not exist\n", ConfigFile)
		return errors.New(str)
	}
	f, err := os.Open(configFile)
	if err != nil {
	}
	defer f.Close()

	configData, err := ioutil.ReadAll(bufio.NewReader(f))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(configData, ServerConfig); err != nil {
		return err
	}
	return nil
}
