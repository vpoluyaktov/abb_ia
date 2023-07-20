package config

import "time"

// singleton
var (
	configInstance *Config
	appVersion, buildDate string
)

type Config struct {
	logFileName        string
	logLevel           string
	useMock            bool
	saveMock           bool
	searchCondition    string
	parrallelDownloads int
	parrallelEncoders  int
	reEncodeFiles      bool
	bitRate            string
	sampleRate         string
}

func Load() {
	config := &Config{}
	config.logFileName = "abb_ia.log"
	config.logLevel = "INFO"
	config.useMock = false
	config.saveMock = false
	config.searchCondition = ""
	config.parrallelDownloads = 5
	config.parrallelEncoders = 5
	config.reEncodeFiles = true
	config.bitRate = "128k"
	config.sampleRate = "44100"

	// read config file here

	configInstance = config
}

func SetLogfileName(fileName string) {
	configInstance.logFileName = fileName
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

func SetParallelDownloads(n int) {
	configInstance.parrallelDownloads = n
}

func ParallelDownloads() int {
	return configInstance.parrallelDownloads
}

func SetParallelEncoders(n int) {
	configInstance.parrallelEncoders = n
}

func ParallelEncoders() int {
	return configInstance.parrallelEncoders
}

func SetReEncodeFiles(b bool) {
	configInstance.reEncodeFiles = b
}

func IsReEncodeFiles() bool {
	return configInstance.reEncodeFiles
}

func SetBitRate(b string) {
	configInstance.bitRate = b
}

func BitRate() string {
	return configInstance.bitRate
}

func SetSampleRate(b string) {
	configInstance.sampleRate = b
}

func SampleRate() string {
	return configInstance.sampleRate
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
