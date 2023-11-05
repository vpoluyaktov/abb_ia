package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/utils"
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
	LogFileName            string
	OutputDir              string
	LogLevel               string
	MaxSearchRows          int
	UseMock                bool
	SaveMock               bool
	SearchCondition        string
	ParrallelDownloads     int
	ParrallelEncoders      int
	ReEncodeFiles          bool
	BitRate                string
	SampleRate             string
	MaxFileSize            string
	CopyToAudiobookshelf   bool
	AudiobookshelfUrl      string
	AudiobookshelfUser     string
	AudiobookshelfPassword string
	AudiobookshelfLibrary  string
	AudiobookshelfDir      string
	ShortenTitles          bool
	Genres                 []string
}

func Load() {
	config := &Config{}

	// default settings
	config.LogFileName = "abb_ia.log"
	config.OutputDir = "output"
	config.LogLevel = "INFO"
	config.MaxSearchRows = 100
	config.UseMock = false
	config.SaveMock = false
	config.SearchCondition = ""
	config.ParrallelDownloads = 5
	config.ParrallelEncoders = 5
	config.ReEncodeFiles = true
	config.BitRate = "128k"
	config.SampleRate = "44100"
	config.MaxFileSize = "100 Mb"
	config.CopyToAudiobookshelf = true
	config.AudiobookshelfUser = "admin"
	config.AudiobookshelfDir = "/mnt/NAS/Audiobooks/Internet Archive"
	config.ShortenTitles = true
	config.Genres = []string{
		"Audiobook",
		"Fiction",
		"Radiodrama",
		"History",
		"Podcast",
		"Nonfiction",
		"Speech",
	}

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

func SetOutputDir(outputDir string) {
	configInstance.OutputDir = outputDir
}

func OutputDir() string {
	return configInstance.OutputDir
}

func SetLogLevel(logLevel string) {
	configInstance.LogLevel = logLevel
}

func LogLevel() string {
	return configInstance.LogLevel
}

func SetMaxSearchRows(r int) {
	configInstance.MaxSearchRows = r
}

func MaxSearchRows() int {
	return configInstance.MaxSearchRows
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
	maxFileSize, err := utils.HumanToBytes(configInstance.MaxFileSize)
	if err != nil {
		logger.Error("Config Loader Can't parse MaxFileSize: " + err.Error() + ". Using default 100 Mb")
		maxFileSize, _ = utils.HumanToBytes("100 Mb")
	}
	return maxFileSize
}

func SetCopyToAudiobookshelf(b bool) {
	configInstance.CopyToAudiobookshelf = b
}

func IsCopyToAudiobookshelf() bool {
	return configInstance.CopyToAudiobookshelf
}

func AudiobookshelfDir() string {
	return configInstance.AudiobookshelfDir
}

func SetAudiobookshelfDir(d string) {
	configInstance.AudiobookshelfDir = d
}

func AudiobookshelfUrl() string {
	return configInstance.AudiobookshelfUrl
}

func SetAudiobookshelfUrl(url string) {
	configInstance.AudiobookshelfUrl = url
}

func AudiobookshelfUser() string {
	return configInstance.AudiobookshelfUser
}

func SetAudiobookshelfUser(u string) {
	configInstance.AudiobookshelfUser = u
}

func AudiobookshelfPassword() string {
	base64 := configInstance.AudiobookshelfPassword
	if base64 == "" {
		return ""
	}
	encrypted, err := utils.DecodeBase64(base64)
	if err != nil {
		logger.Error("Can't decode base64 password: " + err.Error())
		return ""
	}
	decrypted, err := utils.DecryptString(encrypted)
	if err != nil {
		logger.Error("Can't decrypt password: " + err.Error())
		return ""
	}
	return decrypted
}

func SetAudiobookshelfPassword(p string) {

	encrypted, err := utils.EncryptString(p)
	if err != nil {
		logger.Error("Can't encrypt password: " + err.Error())
	}
	base64 := utils.EncodeBase64(encrypted)
	configInstance.AudiobookshelfPassword = base64
}

func AudiobookshelfLibrary() string {
	return configInstance.AudiobookshelfLibrary
}

func SetAudiobookshelfLibrary(l string) {
	configInstance.AudiobookshelfLibrary = l
}

func SetShortenTitles(b bool) {
	configInstance.ShortenTitles = b
}

func IsShortenTitle() bool {
	return configInstance.ShortenTitles
}

func Genres() []string {
	return configInstance.Genres
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
