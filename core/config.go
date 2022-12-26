package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/mlange-42/track/fs"
	"gopkg.in/yaml.v3"
)

// Config for track
type Config struct {
	TextEditor       string        `yaml:"textEditor"`
	MaxBreakDuration time.Duration `yaml:"maxBreakDuration"`
}

// LoadConfig loads the track config, or creates and saves default settings
func LoadConfig() (Config, error) {
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

	conf = Config{
		TextEditor:       editor,
		MaxBreakDuration: 2 * time.Hour,
	}

	err = SaveConfig(conf)
	if err != nil {
		return Config{}, fmt.Errorf("could not save config file: %s", err)
	}

	return conf, nil
}

func tryLoadConfig() (Config, error) {
	file, err := ioutil.ReadFile(fs.ConfigPath())
	if err != nil {
		return Config{}, err
	}

	var conf Config

	if err := yaml.Unmarshal(file, &conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}

// SaveConfig saves the given config to it's default location
func SaveConfig(conf Config) error {
	path := fs.ConfigPath()

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	bytes, err := yaml.Marshal(&conf)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "# Track config\n\n")
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}
