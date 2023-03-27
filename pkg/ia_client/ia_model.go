package ia_client

import (
	"time"
)

type SearchResult struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
		Params struct {
			Query  string `json:"query"`
			Qin    string `json:"qin"`
			Fields string `json:"fields"`
			Wt     string `json:"wt"`
			Rows   string `json:"rows"`
			Start  int    `json:"start"`
		} `json:"params"`
	} `json:"responseHeader"`
	Response struct {
		NumFound int `json:"numFound"`
		Start    int `json:"start"`
		Docs     []struct {
			Collection         []string    `json:"collection"`
			Creator            string      `json:"creator,omitempty"`
			Date               time.Time   `json:"date,omitempty"`
			Description        string      `json:"description,omitempty"`
			Downloads          int         `json:"downloads"`
			Format             []string    `json:"format"`
			Identifier         string      `json:"identifier"`
			Indexflag          []string    `json:"indexflag"`
			ItemSize           int         `json:"item_size"`
			Mediatype          string      `json:"mediatype"`
			Month              int         `json:"month"`
			OaiUpdatedate      []time.Time `json:"oai_updatedate"`
			Publicdate         time.Time   `json:"publicdate"`
			Subject            strArray    `json:"subject,omitempty"`
			Title              string      `json:"title"`
			Week               int         `json:"week"`
			Year               int         `json:"year,omitempty"`
			BackupLocation     string      `json:"backup_location,omitempty"`
			ExternalIdentifier string      `json:"external-identifier,omitempty"`
			Genre              string      `json:"genre,omitempty"`
			Language           string      `json:"language,omitempty"`
			Licenseurl         string      `json:"licenseurl,omitempty"`
			StrippedTags       strArray    `json:"stripped_tags,omitempty"`
		} `json:"docs"`
	} `json:"response"`
}

type GetItemResult struct {
	Server   string `json:"server"`
	Dir      string `json:"dir"`
	Metadata struct {
		Identifier           []string `json:"identifier"`
		Creator              []string `json:"creator"`
		Date                 []string `json:"date"`
		Description          []string `json:"description"`
		GUID                 []string `json:"guid"`
		Mediatype            []string `json:"mediatype"`
		Rssfeed              []string `json:"rssfeed"`
		Scanner              []string `json:"scanner"`
		Sessionid            []string `json:"sessionid"`
		Subject              []string `json:"subject"`
		Title                []string `json:"title"`
		Uploadsoftware       []string `json:"uploadsoftware"`
		Collection           []string `json:"collection"`
		Publicdate           []string `json:"publicdate"`
		Addeddate            []string `json:"addeddate"`
		Curation             []string `json:"curation"`
		AccessRestrictedItem []string `json:"access-restricted-item"`
	} `json:"metadata"`
	Files map[string]struct {
			Source             string `json:"source"`
			Format             string `json:"format"`
			Length             string `json:"length"`
			Mtime              string `json:"mtime"`
			Size               string `json:"size"`
			Md5                string `json:"md5"`
			Crc32              string `json:"crc32"`
			Sha1               string `json:"sha1"`
			Title              string `json:"title"`
			Creator            string `json:"creator"`
			Album              string `json:"album"`
			Artist             string `json:"artist"`
			Genre              string `json:"genre"`
			ExternalIdentifier string `json:"external-identifier"`
			Height             string `json:"height"`
			Width              string `json:"width"`
			Track              string `json:"track"`
			Comment            string `json:"comment"`
	} `json:"files"`
	Misc struct {
		Image           string `json:"image"`
		CollectionTitle string `json:"collection-title"`
	} `json:"misc"`
	Item struct {
		Downloads            int `json:"downloads"`
		Month                int `json:"month"`
		ItemSize             int `json:"item_size"`
		FilesCount           int `json:"files_count"`
		ItemCount            any `json:"item_count"`
		CollectionFilesCount any `json:"collection_files_count"`
		CollectionSize       any `json:"collection_size"`
	} `json:"item"`
}