package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/mlange-42/track/fs"
)

type config struct {
	TextEditor string
}

func loadConfig() (config, error) {
	conf, err := tryLoadConfig()
	if err == nil {
		return conf, nil
	}

	var editor string
	if strings.ToLower(runtime.GOOS) == "windows" {
		editor = "notepad.exe"
	} else {
		editor = "nano"
	}

	conf = config{
		TextEditor: editor,
	}

	err = saveConfig(conf)
	if err != nil {
		return config{}, fmt.Errorf("could not save config file: %s", err)
	}

	return conf, nil
}

func tryLoadConfig() (config, error) {
	file, err := ioutil.ReadFile(fs.ConfigPath())
	if err != nil {
		return config{}, err
	}

	var conf config

	if err := json.Unmarshal(file, &conf); err != nil {
		return config{}, err
	}

	return conf, nil
}

func saveConfig(conf config) error {
	path := fs.ConfigPath()

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(&conf, "", "\t")
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}
