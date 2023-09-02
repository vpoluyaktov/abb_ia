package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// singleton
var (
	configInstance        *Config
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
	viper.SetDefault("LogFileName", "abb_ia.log")
	viper.SetDefault("LogLevel", "INFO")
	viper.SetDefault("UseMock", false)
	viper.SetDefault("SaveMock", false)
	viper.SetDefault("SearchCondition", "")
	viper.SetDefault("ParrallelDownloads", 5)
	viper.SetDefault("ParrallelEncoders", 5)
	viper.SetDefault("ReEncodeFiles", false)
	viper.SetDefault("BitRate", "128k")
	viper.SetDefault("SampleRate", 44100)
	viper.SetDefault("MaxFileSize", 1024*1024)

	viper.SetConfigType("yaml")
	viper.SetConfigFile("./config.yaml")

	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Config file not found. Using default values\n")
		viper.WriteConfig()
	}
	err := viper.Unmarshal(config)
	if err != nil {
		fmt.Printf("Can't parse config file\n")
	}
	configInstance = config
}

func SaveConfig() {
	viper.WriteConfig()
}

func SetLogfileName(fileName string) {
	configInstance.LogFileName = fileName
	viper.Set("logfile", fileName)
	SaveConfig()
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
