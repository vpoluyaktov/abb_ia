package audiobookshelf

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User                 User           `json:"user"`
	UserDefaultLibraryID string         `json:"userDefaultLibraryId"`
	ServerSettings       ServerSettings `json:"serverSettings"`
	Source               string         `json:"Source"`
}

type User struct {
	ID                              string          `json:"id"`
	Username                        string          `json:"username"`
	Type                            string          `json:"type"`
	Token                           string          `json:"token"`
	MediaProgress                   []MediaProgress `json:"mediaProgress"`
	SeriesHideFromContinueListening []interface{}   `json:"seriesHideFromContinueListening"`
	Bookmarks                       []interface{}   `json:"bookmarks"`
	IsActive                        bool            `json:"isActive"`
	IsLocked                        bool            `json:"isLocked"`
	LastSeen                        int64           `json:"lastSeen"`
	CreatedAt                       int64           `json:"createdAt"`
	Permissions                     Permissions     `json:"permissions"`
	LibrariesAccessible             []interface{}   `json:"librariesAccessible"`
	ItemTagsAccessible              []interface{}   `json:"itemTagsAccessible"`
}

type MediaProgress struct {
	ID                        string  `json:"id"`
	LibraryItemID             string  `json:"libraryItemId"`
	EpisodeID                 string  `json:"episodeId"`
	Duration                  float64 `json:"duration"`
	Progress                  float64 `json:"progress"`
	CurrentTime               float64 `json:"currentTime"`
	IsFinished                bool    `json:"isFinished"`
	HideFromContinueListening bool    `json:"hideFromContinueListening"`
	LastUpdate                int64   `json:"lastUpdate"`
	StartedAt                 int64   `json:"startedAt"`
	FinishedAt                *int64  `json:"finishedAt"`
}

type Permissions struct {
	Download              bool `json:"download"`
	Update                bool `json:"update"`
	Delete                bool `json:"delete"`
	Upload                bool `json:"upload"`
	AccessAllLibraries    bool `json:"accessAllLibraries"`
	AccessAllTags         bool `json:"accessAllTags"`
	AccessExplicitContent bool `json:"accessExplicitContent"`
}

type ServerSettings struct {
	ID                                string   `json:"id"`
	ScannerFindCovers                 bool     `json:"scannerFindCovers"`
	ScannerCoverProvider              string   `json:"scannerCoverProvider"`
	ScannerParseSubtitle              bool     `json:"scannerParseSubtitle"`
	ScannerPreferAudioMetadata        bool     `json:"scannerPreferAudioMetadata"`
	ScannerPreferOpfMetadata          bool     `json:"scannerPreferOpfMetadata"`
	ScannerPreferMatchedMetadata      bool     `json:"scannerPreferMatchedMetadata"`
	ScannerDisableWatcher             bool     `json:"scannerDisableWatcher"`
	ScannerPreferOverdriveMediaMarker bool     `json:"scannerPreferOverdriveMediaMarker"`
	ScannerUseSingleThreadedProber    bool     `json:"scannerUseSingleThreadedProber"`
	ScannerMaxThreads                 int      `json:"scannerMaxThreads"`
	ScannerUseTone                    bool     `json:"scannerUseTone"`
	StoreCoverWithItem                bool     `json:"storeCoverWithItem"`
	StoreMetadataWithItem             bool     `json:"storeMetadataWithItem"`
	MetadataFileFormat                string   `json:"metadataFileFormat"`
	RateLimitLoginRequests            int      `json:"rateLimitLoginRequests"`
	RateLimitLoginWindow              int      `json:"rateLimitLoginWindow"`
	BackupSchedule                    string   `json:"backupSchedule"`
	BackupsToKeep                     int      `json:"backupsToKeep"`
	MaxBackupSize                     int      `json:"maxBackupSize"`
	BackupMetadataCovers              bool     `json:"backupMetadataCovers"`
	LoggerDailyLogsToKeep             int      `json:"loggerDailyLogsToKeep"`
	LoggerScannerLogsToKeep           int      `json:"loggerScannerLogsToKeep"`
	HomeBookshelfView                 int      `json:"homeBookshelfView"`
	BookshelfView                     int      `json:"bookshelfView"`
	SortingIgnorePrefix               bool     `json:"sortingIgnorePrefix"`
	SortingPrefixes                   []string `json:"sortingPrefixes"`
	ChromecastEnabled                 bool     `json:"chromecastEnabled"`
	DateFormat                        string   `json:"dateFormat"`
	Language                          string   `json:"language"`
	LogLevel                          int      `json:"logLevel"`
	Version                           string   `json:"version"`
}

type LibrariesResponse struct {
	Libraries []Library `json:"libraries"`
}

type Library struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Folders      []Folder        `json:"folders"`
	DisplayOrder int             `json:"displayOrder"`
	Icon         string          `json:"icon"`
	MediaType    string          `json:"mediaType"`
	Provider     string          `json:"provider"`
	Settings     LibrarySettings `json:"settings"`
	CreatedAt    int64           `json:"createdAt"`
	LastUpdate   int64           `json:"lastUpdate"`
}

type Folder struct {
	ID        string `json:"id"`
	FullPath  string `json:"fullPath"`
	LibraryID string `json:"libraryId"`
	AddedAt   int64  `json:"addedAt,omitempty"`
}

type LibrarySettings struct {
	CoverAspectRatio          float64 `json:"coverAspectRatio"`
	DisableWatcher            bool    `json:"disableWatcher"`
	SkipMatchingMediaWithASIN bool    `json:"skipMatchingMediaWithAsin"`
	SkipMatchingMediaWithISBN bool    `json:"skipMatchingMediaWithIsbn"`
	AutoScanCronExpression    string  `json:"autoScanCronExpression,omitempty"`
}
