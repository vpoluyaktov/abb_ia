package ia_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)


type IAClient struct {
	restyClient *resty.Client 
	
}

type SearchResult struct {
	ResponseHeader ResponseHeader `json:"responseHeader"`
	Response       Response       `json:"response"`
}
type Params struct {
	Query  string `json:"query"`
	Qin    string `json:"qin"`
	Fields string `json:"fields"`
	Wt     string `json:"wt"`
	Rows   string `json:"rows"`
	Start  int    `json:"start"`
}
type ResponseHeader struct {
	Status int    `json:"status"`
	QTime  int    `json:"QTime"`
	Params Params `json:"params"`
}
type Docs struct {
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
	StrippedTags       string      `json:"stripped_tags,omitempty"`
}
type Response struct {
	NumFound int    `json:"numFound"`
	Start    int    `json:"start"`
	Docs     []Docs `json:"docs"`
}

// StrArray string array to be used on JSON UnmarshalJSON
type strArray []string


var (
	// ErrUnsupportedType is returned if the type is not implemented
	ErrUnsupportedType = errors.New("unsupported type")
)

// UnmarshalJSON convert JSON object array of string or
// a string format strings to a golang string array
func (sa *strArray) UnmarshalJSON(data []byte) error {
	var jsonObj interface{}
	err := json.Unmarshal(data, &jsonObj)
	if err != nil {
		return err
	}
	switch obj := jsonObj.(type) {
	case string:
		*sa = strArray([]string{obj})
		return nil
	case []interface{}:
		s := make([]string, 0, len(obj))
		for _, v := range obj {
			value, ok := v.(string)
			if !ok {
				return ErrUnsupportedType
			}
			s = append(s, value)
		}
		*sa = strArray(s)
		return nil
	}
	return ErrUnsupportedType
}


func init() {
}

func New() (* IAClient) {
		var client IAClient
		client.restyClient = resty.New()
		return &client
}

func (client *IAClient) Search(query string) {

	var searchURL = "https://archive.org/advancedsearch.php?q=title:(%s)+AND+mediatype:(audio)&output=json&rows=25&page=1"

	result := &SearchResult{}
	_, err := client.restyClient.R().SetResult(result).Get(fmt.Sprintf(searchURL, query))
	if err != nil {
		panic(err)
	}
  
	fmt.Println(result)
}

