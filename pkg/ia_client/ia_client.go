package ia_client

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
)

const (
	IA_BASE_URL = "https://archive.org"
	MAX_RESULT_ROWS = 25
)

type IAClient struct {
	restyClient *resty.Client 
}

func New() (* IAClient) {
		var client IAClient
		client.restyClient = resty.New()
		return &client
}

func (client *IAClient) SearchByTitle(title string, mediaType string)(*SearchResult) {
	var searchURL = IA_BASE_URL + "/advancedsearch.php?q=title:(%s)+AND+mediatype:(%s)&output=json&rows=%d&page=1"
	result := &SearchResult{}
	_, err := client.restyClient.R().SetResult(result).Get(fmt.Sprintf(searchURL, title, mediaType, MAX_RESULT_ROWS))
	if err != nil {
		panic(err)
	}
	logger.Debug("SearchByTitle response: " + result.String())
	return result
}

func (client *IAClient) SearchByID(itemId string, mediaType string)(*SearchResult) {
	var searchURL = IA_BASE_URL + "/advancedsearch.php?q=identifier:(%s)+AND+mediatype:(%s)&output=json&rows=%d&page=1"
	result := &SearchResult{}
	_, err := client.restyClient.R().SetResult(result).Get(fmt.Sprintf(searchURL, itemId, mediaType, MAX_RESULT_ROWS))
	if err != nil {
		panic(err)
	}
	logger.Debug("SearchByID response: " + result.String())
	return result
}

func (client *IAClient) Search(searchCondition string, mediaType string)(*SearchResult) {
	if strings.Contains(searchCondition, IA_BASE_URL + "/details/") {
		item_id := strings.Split(searchCondition, "/")[4]
		return client.SearchByID(item_id, mediaType)
	} else {
		return client.SearchByTitle(searchCondition, mediaType)
	}
}

func (client *IAClient) GetItemById(itemId string)(*GetItemResult) {
	var getURL = IA_BASE_URL + "/details/%s/?output=json"
	result := &GetItemResult{}
	_, err := client.restyClient.R().SetResult(result).Get(fmt.Sprintf(getURL, itemId))
	if err != nil {
		panic(err)
	}
	logger.Debug("GetItemById response: " + result.String())

	return result
}

func (client *IAClient) DownloadFile() {
	
}

