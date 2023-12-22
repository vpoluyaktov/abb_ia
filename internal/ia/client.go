package ia_client

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"abb_ia/internal/logger"
	"abb_ia/internal/utils"
	"github.com/go-resty/resty/v2"
)

const (
	IA_BASE_URL = "https://archive.org"
	MOCK_DIR    = "mock"
)

type IAClient struct {
	restyClient    *resty.Client
	maxSearchRows  int
	page           int
	loadMockResult bool
	saveMockResult bool
}

func New(maxSearchRows int, useMock bool, saveMock bool) *IAClient {
	client := &IAClient{}
	client.maxSearchRows = maxSearchRows
	client.loadMockResult = useMock
	client.saveMockResult = saveMock

	if client.saveMockResult {
		if err := os.MkdirAll(MOCK_DIR, 0750); err != nil {
			logger.Error("IA Client can't create Mock directory " + MOCK_DIR + ": " + err.Error())
		}
	}
	client.restyClient = resty.New()
	return client
}

func (client *IAClient) Search(author string, title string, mediaType string) *SearchResponse {
	if strings.Contains(title, IA_BASE_URL+"/details/") {
		item_id := strings.Split(title, "/")[4]
		return client.searchByID(item_id, mediaType)
	} else {
		client.page = 1
		return client.searchByTitle(author, title, mediaType)
	}
}

func (client *IAClient) GetNextPage(author string, title string, mediaType string) *SearchResponse {
	if strings.Contains(title, IA_BASE_URL+"/details/") {
		return &SearchResponse{}
	} else {
		client.page += 1
		resp := client.searchByTitle(author, title, mediaType)
		return resp
	}
}

func (client *IAClient) searchByTitle(author string, title string, mediaType string) *SearchResponse {
	mockFile := MOCK_DIR + "/SearchByAuthorAndTitle.json"
	result := &SearchResponse{}
	if client.loadMockResult {
		if err := utils.LoadJson(mockFile, result); err != nil {
			logger.Error("IA Client SearchByAuthorAndTitle() mock load error: " + err.Error())
		}
	} else {
		searchCondition := ""
		if author != "" && title != "" {
			searchCondition = fmt.Sprintf("creator:(%s)+AND+title:(%s)", url.QueryEscape(author), url.QueryEscape(title))
		} else if author != "" {
			searchCondition = fmt.Sprintf("creator:(%s)", url.QueryEscape(author))
		} else if title != "" {
			searchCondition = fmt.Sprintf("title:(%s)", url.QueryEscape(title))
		}
		var searchURL = fmt.Sprintf(IA_BASE_URL+"/advancedsearch.php?q=%s+AND+mediatype:(%s)&output=json&rows=%d&page=%d",
			searchCondition, mediaType, client.maxSearchRows, client.page)
		logger.Debug("IA request: " + searchURL)
		_, err := client.restyClient.R().SetResult(result).Get(searchURL)
		if err != nil {
			logger.Error("IAClient SearchByTitle() error: " + err.Error())
		}
		if client.saveMockResult {
			if err := utils.DumpJson(mockFile, result); err != nil {
				logger.Error("IAClient SearchByTitle() mock save error: " + err.Error())
			}
		}
	}
	// logger.Debug("SearchByTitle response: " + result.String())
	return result
}

func (client *IAClient) searchByID(itemId string, mediaType string) *SearchResponse {
	mockFile := MOCK_DIR + "/SearchByID.json"
	result := &SearchResponse{}
	if client.loadMockResult {
		if err := utils.LoadJson(mockFile, result); err != nil {
			logger.Error("IAClient SearchByID() mock load error: " + err.Error())
		}
	} else {
		var searchURL = fmt.Sprintf(IA_BASE_URL+"/advancedsearch.php?q=identifier:(%s)+AND+mediatype:(%s)&output=json&rows=%d&page=1",
			itemId, mediaType, client.maxSearchRows)
		logger.Debug("IA request: " + searchURL)
		_, err := client.restyClient.R().SetResult(result).Get(searchURL)
		if err != nil {
			logger.Error("IAClient SearchByID() error: " + err.Error())
		}
	}
	if client.saveMockResult {
		if err := utils.DumpJson(mockFile, result); err != nil {
			logger.Error("IAClient SearchByID() mock save error: " + err.Error())
		}
	}
	// logger.Debug("SearchByID response: " + result.String())
	return result
}

func (client *IAClient) GetItemDetails(itemId string) *ItemDetails {
	mockFile := MOCK_DIR + "/GetItemDetails_" + itemId + ".json"
	result := &ItemDetails{}
	if client.loadMockResult {
		delay := time.Duration(rand.Intn(100))
		time.Sleep(3 * delay * time.Millisecond)
		if err := utils.LoadJson(mockFile, result); err != nil {
			logger.Error("IAClient GetItemDetails() mock load error: " + err.Error())
		}
	} else {
		var getURL = fmt.Sprintf(IA_BASE_URL+"/details/%s/?output=json", itemId)
		_, err := client.restyClient.R().SetResult(result).Get(getURL)
		if err != nil {
			logger.Error("IAClient GetItemDetails() error: " + err.Error())
		}
	}
	if client.saveMockResult {
		if err := utils.DumpJson(mockFile, result); err != nil {
			logger.Error("IAClient SearchByID() mock save error: " + err.Error())
		}
	}
	// logger.Debug("GetItemDetails response: " + result.String())

	return result
}

func (client *IAClient) DownloadFile(localDir string, localFile string, iaServer string, iaDir string, iaFile string, fileId int, estimatedSize int64, updateProgress Fn) {

	if client.loadMockResult {
		delay := time.Duration(rand.Intn(10)) // 100
		for percent := 0; percent <= 100; percent++ {
			updateProgress(fileId, iaFile, estimatedSize, int64(float32(estimatedSize)*float32(percent)/100), percent)
			time.Sleep(delay * time.Millisecond)
		}
		return
	}

	iaDir = strings.TrimPrefix(iaDir, "/")
	iaFile = strings.TrimPrefix(iaFile, "/")
	URL := &url.URL{
		Scheme: "https",
		Host:   iaServer,
		Path:   iaDir + "/" + iaFile,
	}
	fileUrl := URL.String()
	localPath := filepath.Join(localDir, localFile)
	tempPath := localPath + ".tmp"

	req, _ := http.NewRequest("GET", fileUrl, nil)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		logger.Fatal("Error while downloading: " + resp.Status)
	}
	defer resp.Body.Close()

	tempDir := filepath.Dir(tempPath)
	if err := os.MkdirAll(tempDir, 0750); err != nil {
		logger.Fatal("Can't create output directory: " + err.Error())
		return
	}
	f, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatal("Can't create temporary file: " + err.Error())
		return
	}
	defer f.Close()

	progressReader := &ProgressReader{
		FileId:   fileId,
		FileName: iaFile,
		Reader:   resp.Body,
		Size:     resp.ContentLength,
		Callback: updateProgress,
	}

	if _, err := io.Copy(f, progressReader); err != nil {
		logger.Fatal("Error while downloading: " + err.Error())
	}

	// fix incorrect ContentLength problem
	if !client.loadMockResult {
		updateProgress(fileId, iaFile, resp.ContentLength, resp.ContentLength, 100)
	}

	os.Rename(tempPath, localPath)
	logger.Debug(iaFile + " downloaded to " + localPath)
}
