package config

import (
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
)

// singleton
var (
	configInstance *Config
)

type Config struct {
	logFileName     string 
	logLevel        string
	useMock         bool
	saveMock        bool
	searchCondition string
}

func Load() {
	config := &Config{}
	config.logFileName = "/tmp/audiobook_creator_IA.log"
	config.logLevel = "DEBUG"
	config.useMock = false
	config.saveMock = false
	config.searchCondition = ""

	// read config file here

	configInstance = config
}

func SetLogfileName(fileName string) {
	configInstance.logFileName = fileName
	logger.SetLogfileName(fileName)
}

func LogFileName() string {
	return configInstance.logFileName
}

func SetLogLevel(logLevel string) {
	configInstance.logLevel = logLevel
}

func LogLevel() string {
	return configInstance.logLevel
}

func UseMock(b bool) {
	configInstance.useMock = b
}

func IsUseMock() bool {
	return configInstance.useMock
}

func SaveMock(b bool) {
	configInstance.saveMock = b
}

func IsSaveMock() bool {
	return configInstance.saveMock
}

func SetSearchCondition(c string) {
	configInstance.searchCondition = c
}

func SearchCondition() string {
	return configInstance.searchCondition
}