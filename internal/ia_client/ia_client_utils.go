package ia_client

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"jaytaylor.com/html2text"
)

// String representation of the GetItemResult struct
func (getItemResult ItemDetails) String() string {
	str, _ := json.Marshal(getItemResult)
	return string(str)
}

// String representation of the SearchResult struct
func (searchResult SearchResponse) String() string {
	str, _ := json.Marshal(searchResult)
	return string(str)
}

// StrArray string array to be used on JSON UnmarshalJSON
type strArray []string

// NumArray int array to be used on JSON UnmarshalJSON
type numArray []float64

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
	if jsonObj == nil {
		*sa = strArray(make([]string, 0))
		return nil
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

// UnmarshalJSON convert JSON object array of int or
// a int format int to a golang int array
func (sa *numArray) UnmarshalJSON(data []byte) error {
	var jsonObj interface{}
	err := json.Unmarshal(data, &jsonObj)
	if err != nil {
		return err
	}
	if jsonObj == nil {
		*sa = numArray(make([]float64, 0))
		return nil
	}	
	switch obj := jsonObj.(type) {
	case float64:
		*sa = numArray([]float64{obj})
		return nil
	case []interface{}:
		i := make([]float64, 0, len(obj))
		for _, v := range obj {
			value, ok := v.(float64)
			if !ok {
				return ErrUnsupportedType
			}
			i = append(i, value)
		}
		*sa = numArray(i)
		return nil
	}
	return ErrUnsupportedType
}

type Fn func(fileId int, fileName string, size int64, pos int64, percent int)

// Progress Reader for file download progress
type ProgressReader struct {
	FileId int
	FileName string
	Reader   io.Reader
	Size     int64
	Pos      int64
	Percent  int
	Callback Fn
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	if err == nil {
		pr.Pos += int64(n)
		pr.Percent = int(float64(pr.Pos) / float64(pr.Size) * 100)
		pr.Callback(pr.FileId, pr.FileName, pr.Size, pr.Pos, pr.Percent)
	}
	return n, err
}

func (client *IAClient) Html2Text(html string) string {
	html = RemoveHtmlTag(html, "<blockquote>")
	html = RemoveHtmlTag(html, "<b>")
	text, err := html2text.FromString(html)
	if err != nil {
		text = "HTML parsing error"
	}
	text = strings.Replace(text, "\u00a0", "\n", -1)

	return text
}

func RemoveHtmlTag(html string, tag string) string {
	html = strings.Replace(html, tag, "", -1)
	html = strings.Replace(html, tag[0:1] + "/" + tag[1:], "", -1)
	return html
}
