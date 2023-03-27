package ia_client

import (
	"encoding/json"
	"errors"
)

// String representation of the GetItemResult struct
func (getItemResult GetItemResult) String() string {
	str, _ := json.Marshal(getItemResult)
	return string(str)
}

// String representation of the SearchResult struct
func (searchResult SearchResult) String() string {
	str, _ := json.Marshal(searchResult)
	return string(str)
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