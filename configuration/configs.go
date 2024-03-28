package configuration

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

const CONFIG_FILENAME string = "config.yaml"

type Cfg struct {
	Log LogConfig
}

type LogConfig struct {
	RotationTime     string `yaml:"log_rotation_time"`
	LogInfoDir       string `yaml:"log_info_dir"`
	LogInfoFileName  string `yaml:"log_info_name"`
	LogDebugDir      string `yaml:"log_debug_dir"`
	LogDebugFileName string `yaml:"log_debug_name"`
}

var config *Cfg
var once sync.Once

// Singleton pattern
func GetConfigs() *Cfg {
	once.Do(func() {
		var err error
		if config, err = newConfigs(); err != nil {
			panic(err)
		}
	})
	return config
}

func newConfigs() (*Cfg, error) {
	cfg := readConfEnv()
	if cfg != nil {
		return cfg, nil
	}
	filename := os.Getenv("LOG_CONFIG_NAME")
	if len(filename) == 0 {
		filename = CONFIG_FILENAME
	}
	return readConf(filename)
}

func readConfEnv() *Cfg {
	rotationTime := os.Getenv("LOG_ROTATION_TIME")
	logInfoDir := os.Getenv("LOG_INFO_DIR")
	logInfoFileName := os.Getenv("LOG_INFO_NAME")
	logDebugDir := os.Getenv("LOG_DEBUG_DIR")
	logDebugFileName := os.Getenv("LOG_DEBUG_NAME")
	if len(rotationTime) == 0 || len(logInfoDir) == 0 || len(logInfoFileName) == 0 || len(logDebugDir) == 0 || len(logDebugFileName) == 0 {
		return nil
	}
	return &Cfg{
		LogConfig{
			rotationTime,
			logInfoDir,
			logInfoFileName,
			logDebugDir,
			logDebugFileName,
		},
	}
}

func readConf(filename string) (*Cfg, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Cfg{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	return c, err
}
