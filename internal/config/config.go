package config

// singleton
var (
	configInstance *Config
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
}

func Load() {
	config := &Config{}
	config.logFileName = "/tmp/abb_ia.log"
	config.logLevel = "INFO"
	config.useMock = false
	config.saveMock = false
	config.searchCondition = ""
	config.parrallelDownloads = 10
	config.parrallelEncoders = 1
	config.reEncodeFiles = true

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

func SaveParallelDownloads(n int) {
	configInstance.parrallelDownloads = n
}

func GetParallelDownloads() int {
	return configInstance.parrallelDownloads
}

func SaveParallelEncoders(n int) {
	configInstance.parrallelEncoders = n
}

func GetParallelEncoders() int {
	return configInstance.parrallelEncoders
}

func SaveReEncodeFiles(b bool) {
	configInstance.reEncodeFiles = b
}

func IsReEncodeFiles() bool {
	return configInstance.reEncodeFiles
}
