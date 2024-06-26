package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"abb_ia/internal/logger"
	"abb_ia/internal/utils"

	"gopkg.in/yaml.v3"
)

// singleton
var (
	configInstance *Config
)

// global vars
var (
	configFile            = "abb_ia.config.yaml"
	appVersion, buildDate string
	repoOwner             string = "vpoluyaktov"
	repoName              string = "abb_ia"
)

// Fields of this stuct should to be private but I have to make them public because yaml.Marshal/Unmarshal can't work with private fields
type Config struct {
	DefaultAuthor          string        `yaml:"DefaultAuthor"`
	DefaultTitle           string        `yaml:"DefaultTitle"`
	SortBy                 string        `yaml:"SortBy"`
	SortOrder              string        `yaml:"SortOrder"`
	RowsPerPage            int           `yaml:"RowsPerPage"`
	LogFileName            string        `yaml:"LogFileName"`
	OutputDir              string        `yaml:"Outputdir"`
	CopyToOutputDir        bool          `yaml:"CopyToOutputDir"`
	TmpDir                 string        `yaml:"TmpDir"`
	LogLevel               string        `yaml:"LogLevel"`
	UseMock                bool          `yaml:"UseMock"`
	SaveMock               bool          `yaml:"SaveMock"`
	ConcurrentDownloaders  int           `yaml:"ConcurrentDownloaders"`
	ConcurrentEncoders     int           `yaml:"ConcurrentEncoders"`
	ReEncodeFiles          bool          `yaml:"ReEncodeFiles"`
	BasePortNumber         int           `yaml:"BasePortNumber"`
	BitRateKbs             int           `yaml:"BitRateKbs"`
	SampleRateHz           int           `yaml:"SampleRateHz"`
	MaxFileSizeMb          int           `yaml:"MaxFileSizeMb"`
	UploadToAudiobookshef  bool          `yaml:"UploadToAudiobookshelf"`
	ScanAudiobookshef      bool          `yaml:"ScanAudiobookshelf"`
	AudiobookshelfUrl      string        `yaml:"AudiobookshelfUrl"`
	AudiobookshelfUser     string        `yaml:"AudiobookshelfUser"`
	AudiobookshelfPassword string        `yaml:"AudiobookshelfPassword"`
	AudiobookshelfLibrary  string        `yaml:"AudiobookshelfLibrary"`
	ShortenTitles          bool          `yaml:"ShortenTitles"`
	ShortenPairs           []ShortenPair `yaml:"ShortenPairs"`
	Genres                 []string      `yaml:"Genres"`
}

type ShortenPair struct {
	Search  string `yaml:"Search"`
	Replace string `yaml:"Replace"`
}

func Instance() *Config {
	if configInstance == nil {
		configInstance = &Config{}
	}
	return configInstance
}

func Load() {
	config := &Config{}

	// default settings
	config.LogFileName = "abb_ia.log"
	config.TmpDir = "tmp"
	config.CopyToOutputDir = true
	config.OutputDir = "output"
	config.LogLevel = "INFO"
	config.RowsPerPage = 25
	config.UseMock = false
	config.SaveMock = false
	config.DefaultAuthor = "Old Time Radio Researchers Group"
	config.DefaultTitle = "Single Episodes"
	config.SortBy = "Date"
	config.SortOrder = "Descending"
	config.ConcurrentDownloaders = 5
	config.ConcurrentEncoders = 5
	config.ReEncodeFiles = true
	config.BasePortNumber = 31000
	config.BitRateKbs = 128
	config.SampleRateHz = 44100
	config.MaxFileSizeMb = 250
	config.UploadToAudiobookshef = false
	config.ScanAudiobookshef = false
	config.AudiobookshelfUser = "admin"
	config.AudiobookshelfPassword = ""
	config.AudiobookshelfLibrary = "Internet Archive"
	config.ShortenTitles = true
	config.ShortenPairs = []ShortenPair{
		{"Old Time Radio Researchers Group", "OTRR"},
		{" - Single Episodes", ""},
	}
	config.Genres = []string{
		"Audiobook",
		"Fiction",
		"Radiodrama",
		"History",
		"Podcast",
		"Nonfiction",
		"Education",
		"News",
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

func (c *Config) SetLogfileName(fileName string) {
	c.LogFileName = fileName
}

func (c *Config) GetLogFileName() string {
	return c.LogFileName
}

func (c *Config) SetTmpDir(tmpDir string) {
	c.TmpDir = tmpDir
}

func (c *Config) GetTmpDir() string {
	return c.TmpDir
}

func (c *Config) GetOutputDir() string {
	return c.OutputDir
}

func (c *Config) SetOutputdDir(d string) {
	c.OutputDir = d
}

func (c *Config) SetCopyToOutputDir(b bool) {
	c.CopyToOutputDir = b
}

func (c *Config) IsCopyToOutputDir() bool {
	return c.CopyToOutputDir
}

func (c *Config) SetLogLevel(logLevel string) {
	c.LogLevel = logLevel
}

func (c *Config) GetLogLevel() string {
	return c.LogLevel
}

func (c *Config) SetRowsPerPage(r int) {
	c.RowsPerPage = r
}

func (c *Config) GetRowsPerPage() int {
	return c.RowsPerPage
}

func (c *Config) SetUseMock(b bool) {
	c.UseMock = b
}

func (c *Config) IsUseMock() bool {
	return c.UseMock
}

func (c *Config) SetSaveMock(b bool) {
	c.SaveMock = b
}

func (c *Config) IsSaveMock() bool {
	return c.SaveMock
}

func (c *Config) SetDefaultAuthor(s string) {
	c.DefaultAuthor = s
}

func (c *Config) GetDefaultAuthor() string {
	return c.DefaultAuthor
}

func (c *Config) SetDefaultTitle(s string) {
	c.DefaultTitle = s
}

func (c *Config) GetDefaultTitle() string {
	return c.DefaultTitle
}

func (c *Config) SetSortBy(s string) {
	c.SortBy = s
}

func (c *Config) GetSortBy() string {
	return c.SortBy
}

func (c *Config) GetSortByOptions() []string {
	return []string{"Creator", "Title", "Date", "Size"}
}

func (c *Config) SetSortOrder(s string) {
	c.SortOrder = s
}

func (c *Config) GetSortOrder() string {
	return c.SortOrder
}

func (c *Config) GetSortOrderOptions() []string {
	return []string{"Ascending", "Descending"}
}

func (c *Config) SetConcurrentDownloaders(n int) {
	c.ConcurrentDownloaders = n
}

func (c *Config) GetConcurrentDownloaders() int {
	return c.ConcurrentDownloaders
}

func (c *Config) SetConcurrentEncoders(n int) {
	c.ConcurrentEncoders = n
}

func (c *Config) GetConcurrentEncoders() int {
	return c.ConcurrentEncoders
}

func (c *Config) SetReEncodeFiles(b bool) {
	c.ReEncodeFiles = b
}

func (c *Config) IsReEncodeFiles() bool {
	return c.ReEncodeFiles
}

func (c *Config) SetBasePortNumber(port int) {
	c.BasePortNumber = port
}

func (c *Config) GetBasePortNumber() int {
	return c.BasePortNumber
}

func (c *Config) SetBitRate(b int) {
	c.BitRateKbs = b
}

func (c *Config) GetBitRate() int {
	return c.BitRateKbs
}

func (c *Config) SetSampleRate(b int) {
	c.SampleRateHz = b
}

func (c *Config) GetSampleRate() int {
	return c.SampleRateHz
}

func (c *Config) SetMaxFileSizeMb(s int) {
	c.MaxFileSizeMb = s
}

func (c *Config) GetMaxFileSizeMb() int {
	return c.MaxFileSizeMb
}

func (c *Config) SetUploadToAudiobookshelf(b bool) {
	c.UploadToAudiobookshef = b
}

func (c *Config) IsUploadToAudiobookshef() bool {
	return c.UploadToAudiobookshef
}

func (c *Config) SetScanAudiobookshelf(b bool) {
	c.ScanAudiobookshef = b
}

func (c *Config) IsScanAudiobookshef() bool {
	return c.ScanAudiobookshef
}

func (c *Config) GetAudiobookshelfUrl() string {
	return c.AudiobookshelfUrl
}

func (c *Config) SetAudiobookshelfUrl(url string) {
	c.AudiobookshelfUrl = url
}

func (c *Config) GetAudiobookshelfUser() string {
	return c.AudiobookshelfUser
}

func (c *Config) SetAudiobookshelfUser(u string) {
	c.AudiobookshelfUser = u
}

func (c *Config) GetAudiobookshelfPassword() string {
	base64 := c.AudiobookshelfPassword
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

func (c *Config) SetAudiobookshelfPassword(p string) {

	encrypted, err := utils.EncryptString(p)
	if err != nil {
		logger.Error("Can't encrypt password: " + err.Error())
	}
	base64 := utils.EncodeBase64(encrypted)
	c.AudiobookshelfPassword = base64
}

func (c *Config) GetAudiobookshelfLibrary() string {
	return c.AudiobookshelfLibrary
}

func (c *Config) SetAudiobookshelfLibrary(l string) {
	c.AudiobookshelfLibrary = l
}

func (c *Config) SetShortenTitles(b bool) {
	c.ShortenTitles = b
}

func (c *Config) IsShortenTitle() bool {
	return c.ShortenTitles
}

func (c *Config) GetGenres() []string {
	return c.Genres
}

func (c *Config) AppVersion() string {
	if appVersion == "" {
		appVersion = "0.0.0"
	}
	return appVersion
}

func (c *Config) GetRepoOwner() string {
	return repoOwner
}

func (c *Config) GetRepoName() string {
	return repoName
}

func (c *Config) GetBuildDate() string {
	// 2023-07-20T14:45:12Z
	fmt := "01/02/2006"
	bd, err := time.Parse(time.RFC3339, buildDate)
	if buildDate != "" && err != nil {
		return bd.Format(fmt)
	} else {
		return time.Now().Format(fmt)
	}
}

func (c *Config) GetCopy() Config {
	return *configInstance
}
