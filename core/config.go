package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mlange-42/track/fs"
	"gopkg.in/yaml.v3"
)

const defaultWorkspace = "default"

var (
	// ErrNoConfig is an error for no config file available
	ErrNoConfig = errors.New("no config file")
)

// Config for track
type Config struct {
	Workspace        string        `yaml:"workspace"`
	TextEditor       string        `yaml:"textEditor"`
	MaxBreakDuration time.Duration `yaml:"maxBreakDuration"`
	EmptyCell        string        `yaml:"emptyCell"`
	PauseCell        string        `yaml:"pauseCell"`
}

// LoadConfig loads the track config, or creates and saves default settings
func LoadConfig() (Config, error) {
	conf, err := tryLoadConfig()
	if err == nil {
		return conf, nil
	}
	if !errors.Is(err, ErrNoConfig) {
		return conf, err
	}

	var editor string
	if strings.ToLower(runtime.GOOS) == "windows" {
		editor = "notepad.exe"
	} else {
		editor = "nano"
	}

	conf = Config{
		Workspace:        defaultWorkspace,
		TextEditor:       editor,
		MaxBreakDuration: 2 * time.Hour,
		EmptyCell:        ".",
		PauseCell:        "-",
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
		return Config{}, ErrNoConfig
	}

	var conf Config

	if err := yaml.Unmarshal(file, &conf); err != nil {
		return Config{}, err
	}

	if err = CheckConfig(&conf); err != nil {
		return conf, err
	}

	return conf, nil
}

// SaveConfig saves the given config to it's default location
func SaveConfig(conf Config) error {
	if err := CheckConfig(&conf); err != nil {
		return err
	}

	path := fs.ConfigPath()

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	bytes, err := yaml.Marshal(&conf)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "%s Track config\n\n", YamlCommentPrefix)
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}

// CheckConfig checks a config for consistency
func CheckConfig(conf *Config) error {
	if utf8.RuneCountInString(conf.EmptyCell) != 1 {
		return fmt.Errorf("config entry EmptyCell must be a string of length 1. Got '%s'", conf.EmptyCell)
	}
	if utf8.RuneCountInString(conf.PauseCell) != 1 {
		return fmt.Errorf("config entry PauseCell must be a string of length 1. Got '%s'", conf.PauseCell)
	}
	return nil
}
