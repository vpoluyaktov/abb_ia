package ia_client

import (
	"io"
	"os"
	"fmt"
	"strings"
	"net/http"

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

func (client *IAClient) Search(searchCondition string, mediaType string)(*SearchResult) {
	if strings.Contains(searchCondition, IA_BASE_URL + "/details/") {
		item_id := strings.Split(searchCondition, "/")[4]
		return client.searchByID(item_id, mediaType)
	} else {
		return client.searchByTitle(searchCondition, mediaType)
	}
}

func (client *IAClient) searchByTitle(title string, mediaType string)(*SearchResult) {
	var searchURL = IA_BASE_URL + "/advancedsearch.php?q=title:(%s)+AND+mediatype:(%s)&output=json&rows=%d&page=1"
	result := &SearchResult{}
	_, err := client.restyClient.R().SetResult(result).Get(fmt.Sprintf(searchURL, title, mediaType, MAX_RESULT_ROWS))
	if err != nil {
		panic(err)
	}
	logger.Debug("SearchByTitle response: " + result.String())
	return result
}

func (client *IAClient) searchByID(itemId string, mediaType string)(*SearchResult) {
	var searchURL = IA_BASE_URL + "/advancedsearch.php?q=identifier:(%s)+AND+mediatype:(%s)&output=json&rows=%d&page=1"
	result := &SearchResult{}
	_, err := client.restyClient.R().SetResult(result).Get(fmt.Sprintf(searchURL, itemId, mediaType, MAX_RESULT_ROWS))
	if err != nil {
		panic(err)
	}
	logger.Debug("SearchByID response: " + result.String())
	return result
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


func (client *IAClient) DownloadFile(outputDir string, server string, dir string, file string, updateProgress Fn) {
	dir = strings.TrimPrefix(dir, "/")
	file = strings.TrimPrefix(file, "/")
	fileUrl := fmt.Sprintf("https://%s/%s/%s", server, dir, file)
	outPath := fmt.Sprintf("%s/%s/%s", outputDir, dir, file)
	tempPath := outPath + ".tmp"

	req, _ := http.NewRequest("GET", fileUrl, nil)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
			logger.Fatal("Error while downloading: " + resp.Status)
	}
	defer resp.Body.Close()

	if err := os.MkdirAll(fmt.Sprintf("%s/%s", outputDir, dir), 0750); err != nil {
		logger.Fatal("Error while creating output directory: " + err.Error())
	}
	f, _ := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	progressReader := &ProgressReader{
			Reader: resp.Body,
			Size:   resp.ContentLength,
			Callback: updateProgress,
	}

	if _, err := io.Copy(f, progressReader); err != nil {
			logger.Fatal("Error while downloading: " + err.Error())
	}
	os.Rename(tempPath, outPath)
	logger.Debug(file + " downloaded to " + outPath)
}

