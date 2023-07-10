package ia_client

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/utils"
)

const (
	IA_BASE_URL     = "https://archive.org"
	MAX_RESULT_ROWS = 25
	MOCK_DIR        = "mock"
)

type IAClient struct {
	restyClient    *resty.Client
	loadMockResult bool
	saveMockResult bool
}

func New(useMock bool, saveMock bool) *IAClient {
	client := &IAClient{}
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

func (client *IAClient) Search(searchCondition string, mediaType string) *SearchResponse {
	if strings.Contains(searchCondition, IA_BASE_URL+"/details/") {
		item_id := strings.Split(searchCondition, "/")[4]
		return client.searchByID(item_id, mediaType)
	} else {
		return client.searchByTitle(searchCondition, mediaType)
	}
}

func (client *IAClient) searchByTitle(title string, mediaType string) *SearchResponse {
	mockFile := MOCK_DIR + "/SearchByTitle.json"
	result := &SearchResponse{}
	if client.loadMockResult {
		if err := utils.LoadJson(mockFile, result); err != nil {
			logger.Error("IA Client SearchByTitle() mock load error: " + err.Error())
		}
	} else {
		var searchURL = IA_BASE_URL + "/advancedsearch.php?q=title:(%s)+AND+mediatype:(%s)&output=json&rows=%d&page=1"
		_, err := client.restyClient.R().SetResult(result).Get(fmt.Sprintf(searchURL, title, mediaType, MAX_RESULT_ROWS))
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
		var searchURL = IA_BASE_URL + "/advancedsearch.php?q=identifier:(%s)+AND+mediatype:(%s)&output=json&rows=%d&page=1"
		req := fmt.Sprintf(searchURL, itemId, mediaType, MAX_RESULT_ROWS)
		logger.Debug("req: " + req)
		_, err := client.restyClient.R().SetResult(result).Get(fmt.Sprintf(searchURL, itemId, mediaType, MAX_RESULT_ROWS))
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
		var getURL = IA_BASE_URL + "/details/%s/?output=json"
		_, err := client.restyClient.R().SetResult(result).Get(fmt.Sprintf(getURL, itemId))
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

func (client *IAClient) DownloadFile(outputDir string, server string, dir string, fileName string, fileId int, estimatedSize int64, updateProgress Fn) {

	if client.loadMockResult {
		delay := time.Duration(rand.Intn(100))
		for percent := 0; percent <= 100; percent++ {
			updateProgress(fileId, fileName, estimatedSize, int64(float32(estimatedSize)*float32(percent)/100), percent)
			time.Sleep(delay * time.Millisecond)
		}
		return
	}

	dir = strings.TrimPrefix(dir, "/")
	fileName = strings.TrimPrefix(fileName, "/")
	fileUrl := fmt.Sprintf("https://%s/%s/%s", server, dir, fileName)
	outPath := fmt.Sprintf("%s/%s/%s", outputDir, dir, fileName)
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
		FileId:   fileId,
		FileName: fileName,
		Reader:   resp.Body,
		Size:     resp.ContentLength,
		Callback: updateProgress,
	}

	if _, err := io.Copy(f, progressReader); err != nil {
		logger.Fatal("Error while downloading: " + err.Error())
	}

	// fix incorrect ContentLength problem
	if !client.loadMockResult {
		updateProgress(fileId, fileName, resp.ContentLength, resp.ContentLength, 100)
	}

	os.Rename(tempPath, outPath)
	logger.Debug(fileName + " downloaded to " + outPath)
}
