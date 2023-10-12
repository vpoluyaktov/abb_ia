package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

// singleton
var (
	configInstance *Config
)

// global vars
var (
	configFile            = "config.yaml"
	appVersion, buildDate string
)

type Config struct {
	LogFileName        string
	LogLevel           string
	UseMock            bool
	SaveMock           bool
	SearchCondition    string
	ParrallelDownloads int
	ParrallelEncoders  int
	ReEncodeFiles      bool
	BitRate            string
	SampleRate         string
	MaxFileSize        int64
}

func Load() {
	config := &Config{}

	// default settings
	config.LogFileName = "abb_ia.log"
	config.LogLevel = "INFO"
	config.UseMock = false
	config.SaveMock = false
	config.SearchCondition = ""
	config.ParrallelDownloads = 5
	config.ParrallelEncoders = 5
	config.ReEncodeFiles = false
	config.BitRate = "128k"
	config.SampleRate = "44100"
	config.MaxFileSize = 1024 * 1024 * 10

	fmt.Printf("Using config: %s\n", configFile)
	if ReadConfig(config) != nil {
		fmt.Printf("Can read config file. Creating new one\n")
		SaveConfig(config)
	}
	configInstance = config
}

func ReadConfig(c *Config) error {
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	} else {
		err = yaml.Unmarshal(buf, c)
		if err != nil {
			fmt.Printf("Can parse config file. Creating new one\n")
			return err
		}
	}
	return nil
}

func SaveConfig(c *Config) error {
	yaml, err := yaml.Marshal(c)
	if err != nil {
		fmt.Printf("Can not marshal config structure: %s\n", err.Error())
		return err
	} else {
		err = ioutil.WriteFile(configFile, yaml, 0644)
		if err != nil {
			fmt.Printf("Can not write config file: %s\n", err.Error())
		}
		configInstance = c
		return nil
	}
}

func SetLogfileName(fileName string) {
	configInstance.LogFileName = fileName
}

func LogFileName() string {
	return configInstance.LogFileName
}

func SetLogLevel(logLevel string) {
	configInstance.LogLevel = logLevel
}

func LogLevel() string {
	return configInstance.LogLevel
}

func UseMock(b bool) {
	configInstance.UseMock = b
}

func IsUseMock() bool {
	return configInstance.UseMock
}

func SaveMock(b bool) {
	configInstance.SaveMock = b
}

func IsSaveMock() bool {
	return configInstance.SaveMock
}

func SetSearchCondition(c string) {
	configInstance.SearchCondition = c
}

func SearchCondition() string {
	return configInstance.SearchCondition
}

func SetParallelDownloads(n int) {
	configInstance.ParrallelDownloads = n
}

func ParallelDownloads() int {
	return configInstance.ParrallelDownloads
}

func SetParallelEncoders(n int) {
	configInstance.ParrallelEncoders = n
}

func ParallelEncoders() int {
	return configInstance.ParrallelEncoders
}

func SetReEncodeFiles(b bool) {
	configInstance.ReEncodeFiles = b
}

func IsReEncodeFiles() bool {
	return configInstance.ReEncodeFiles
}

func SetBitRate(b string) {
	configInstance.BitRate = b
}

func BitRate() string {
	return configInstance.BitRate
}

func SetSampleRate(b string) {
	configInstance.SampleRate = b
}

func SampleRate() string {
	return configInstance.SampleRate
}

func MaxFileSize() int64 {
	return configInstance.MaxFileSize
}

func AppVersion() string {
	if appVersion == "" {
		appVersion = "0.0.0"
	}
	return appVersion
}

func BuildDate() string {
	// 2023-07-20T14:45:12Z
	fmt := "01/02/2006"
	bd, err := time.Parse(time.RFC3339, buildDate)
	if buildDate != "" && err != nil {
		return bd.Format(fmt)
	} else {
		return time.Now().Format(fmt)
	}
}

func GetCopy() Config {
	return *configInstance
}
