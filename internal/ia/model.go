package ia_client

import (
	"time"
)

type SearchResponse struct {
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
			AvgRating          float64     `json:"avg_rating"`
			Btih               strArray    `json:"btih"`
			Collection         []string    `json:"collection"`
			Creator            strArray    `json:"creator,omitempty"`
			Date               time.Time   `json:"date,omitempty"`
			Description        string      `json:"description,omitempty"`
			Downloads          int         `json:"downloads"`
			Format             strArray    `json:"format"`
			Identifier         string      `json:"identifier"`
			Indexflag          []string    `json:"indexflag"`
			ItemSize           int         `json:"item_size"`
			Mediatype          string      `json:"mediatype"`
			Month              int         `json:"month"`
			OaiUpdatedate      []time.Time `json:"oai_updatedate"`
			Publicdate         time.Time   `json:"publicdate"`
			Reviewdate         time.Time   `json:"reviewdate"`
			Subject            strArray    `json:"subject,omitempty"`
			Title              string      `json:"title"`
			Week               int         `json:"week"`
			Year               numArray    `json:"year,omitempty"`
			BackupLocation     string      `json:"backup_location,omitempty"`
			ExternalIdentifier strArray    `json:"external-identifier,omitempty"`
			Genre              strArray    `json:"genre,omitempty"`
			Language           string      `json:"language,omitempty"`
			Licenseurl         string      `json:"licenseurl,omitempty"`
			StrippedTags       strArray    `json:"stripped_tags,omitempty"`
		} `json:"docs"`
	} `json:"response"`
}

type ItemDetails struct {
	Server   string `json:"server"`
	Dir      string `json:"dir"`
	Metadata struct {
		Identifier             []string `json:"identifier"`
		Creator                []string `json:"creator"`
		Artist                 []string `json:"artist"`
		Date                   []string `json:"date"`
		Description            []string `json:"description"`
		GUID                   []string `json:"guid"`
		Mediatype              []string `json:"mediatype"`
		Rssfeed                []string `json:"rssfeed"`
		Scanner                []string `json:"scanner"`
		Sessionid              []string `json:"sessionid"`
		Subject                []string `json:"subject"`
		Title                  []string `json:"title"`
		Uploadsoftware         []string `json:"uploadsoftware"`
		Collection             []string `json:"collection"`
		Licenseurl             []string `json:"licenseurl"`
		Notes                  []string `json:"notes"`
		Publicdate             []string `json:"publicdate"`
		Addeddate              []string `json:"addeddate"`
		Curation               []string `json:"curation"`
		BackupLocation         []string `json:"backup_location"`
		AccessRestrictedItem   []string `json:"access-restricted-item"`
		ExternalMetadataUpdate []string `json:"external_metadata_update"`
		Reviews                struct {
			Info    map[string]string `json:"info"`
			Reviews []struct {
				ReviewBody       string `json:"reviewbody"`
				ReviewTitle      string `json:"reviewtitle"`
				Reviewer         string `json:"reviewer"`
				ReviewerItemname string `json:"reviewer_itemname"`
				ReviewDate       string `json:"reviewdate"`
				Stars            string `json:"stars"`
			} `json:"reviews"`
		} `json:"reviews"`
	} `json:"metadata"`
	Files map[string]struct {
		Source             string   `json:"source"`
		Format             string   `json:"format"`
		Original           string   `json:"original"`
		Length             string   `json:"length"`
		Mtime              string   `json:"mtime"`
		Size               string   `json:"size"`
		Md5                string   `json:"md5"`
		Crc32              string   `json:"crc32"`
		Sha1               string   `json:"sha1"`
		Title              string   `json:"title"`
		Creator            strArray `json:"creator"`
		Album              string   `json:"album"`
		Artist             string   `json:"artist"`
		Genre              string   `json:"genre"`
		ExternalIdentifier strArray `json:"external-identifier"`
		Height             string   `json:"height"`
		Width              string   `json:"width"`
		Track              string   `json:"track"`
		Comment            string   `json:"comment"`
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
