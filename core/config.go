package core

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

const defaultWorkspace = "default"

var (
	// ErrNoConfig is an error for no config file available
	ErrNoConfig = errors.New("no config file")
)

// Config for Track
type Config struct {
	// The current workspace
	Workspace string `yaml:"workspace"`
	// The text editor for editing resources
	TextEditor string `yaml:"textEditor"`
	// Maximum duration of breaks between records of the same project to consider it as a pause
	MaxBreakDuration time.Duration `yaml:"maxBreakDuration"`
	// Character for empty cells in day and week reports
	EmptyCell string `yaml:"emptyCell"`
	// Character for record cells in day and week reports
	RecordCell string `yaml:"recordCell"`
	// Character for pause cells in day and week reports
	PauseCell string `yaml:"pauseCell"`
}

// defaultConfig creates a Config with default values
func defaultConfig() Config {
	var editor string
	if strings.ToLower(runtime.GOOS) == "windows" {
		editor = "notepad.exe"
	} else {
		editor = "nano"
	}

	return Config{
		Workspace:        defaultWorkspace,
		TextEditor:       editor,
		MaxBreakDuration: 2 * time.Hour,
		EmptyCell:        ".",
		RecordCell:       ":",
		PauseCell:        "-",
	}
}

// LoadConfig loads the track config, or creates and saves default settings
// if it does not exist.
func LoadConfig(path string) (Config, error) {
	conf, err := tryLoadConfig(path)
	if err == nil {
		return conf, nil
	}
	if !errors.Is(err, ErrNoConfig) {
		return conf, err
	}

	conf = defaultConfig()

	err = conf.Save(path)
	if err != nil {
		return Config{}, fmt.Errorf("could not save config file: %s", err)
	}

	return conf, nil
}

func tryLoadConfig(path string) (Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return Config{}, ErrNoConfig
	}

	var conf Config

	if err := yaml.Unmarshal(file, &conf); err != nil {
		return Config{}, err
	}

	if err = conf.Check(); err != nil {
		return conf, err
	}

	return conf, nil
}

// Save saves the given Config to it's default location
func (conf *Config) Save(path string) error {
	if err := conf.Check(); err != nil {
		return err
	}

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

// Check checks the config for consistency
func (conf *Config) Check() error {
	versionHint := "In case you recently updated track, try to delete file %USER%/.track/config.yml"
	if utf8.RuneCountInString(conf.EmptyCell) != 1 {
		return fmt.Errorf("config entry EmptyCell must be a string of length 1. Got '%s'.\n%s", conf.EmptyCell, versionHint)
	}
	if utf8.RuneCountInString(conf.RecordCell) != 1 {
		return fmt.Errorf("config entry RecordCell must be a string of length 1. Got '%s'.\n%s", conf.RecordCell, versionHint)
	}
	if utf8.RuneCountInString(conf.PauseCell) != 1 {
		return fmt.Errorf("config entry PauseCell must be a string of length 1. Got '%s'.\n%s", conf.PauseCell, versionHint)
	}
	return nil
}
