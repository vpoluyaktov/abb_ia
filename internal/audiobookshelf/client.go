package audiobookshelf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/vpoluyaktov/abb_ia/internal/dto"
)

type AudiobookshelfClient struct {
	url           string
	userName      string
	password      string
	loginResponse *LoginResponse
}

func NewClient(url string) *AudiobookshelfClient {
	c := &AudiobookshelfClient{
		url: url,
	}
	return c
}

// Call the Audiobookshelf API login method.
func (c *AudiobookshelfClient) Login(userName string, password string) error {

	c.userName = userName
	c.password = password

	requestBody := LoginRequest{
		Username: c.userName,
		Password: c.password,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal login response body: %v", err)
	}

	resp, err := http.Post(c.url+"/login", "application/json", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return fmt.Errorf("failed to make login API call: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login API call returned status code: %d", resp.StatusCode)
	}

	var loginResp LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		return fmt.Errorf("failed to decode login response: %v", err)
	}
	c.loginResponse = &loginResp

	return nil
}

// Call the Audiobookshelf API libraries method
func (c *AudiobookshelfClient) GetLibraries() ([]Library, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", c.url+"/api/libraries", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+c.loginResponse.User.Token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var response LibrariesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %v", err)
	}

	return response.Libraries, nil
}

func (c *AudiobookshelfClient) GetLibraryId(libraries []Library, libraryName string) (libraryID string, err error) {
	for _, library := range libraries {
		if library.Name == libraryName {
			return library.ID, nil
		}
	}
	return "", fmt.Errorf("no library with name '%s' found", libraryName)
}

func (c *AudiobookshelfClient) GetFolders(libraries []Library, libraryName string) (folders []Folder, err error) {
	for _, library := range libraries {
		if library.Name == libraryName {
			return library.Folders, nil
		}
	}
	return nil, fmt.Errorf("no library with name '%s' found", libraryName)
}


// Call the Audiobookshelf API for a library Scan
func (c *AudiobookshelfClient) ScanLibrary(libraryID string) error {
	url := c.url + "/api/libraries/" + libraryID + "/scan"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.loginResponse.User.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("an admin user is required to start a scan")
	} else if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("the user cannot access the library or no library with the provided ID exists")
	} else {
		return nil
	}
}

// Upload audiobook to The Audibookshelf server
func (c *AudiobookshelfClient) UploadBook(ab *dto.Audiobook, libraryID string, folderID string, callback Fn ) error {
	// Open each file for upload
	var filesList []*os.File
	for _, part := range ab.Parts {
		f, err := os.Open(part.M4BFile)
		if err != nil {
			return err
		}
		defer f.Close()
		filesList = append(filesList, f)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add metadata fields
	_ = writer.WriteField("title", ab.Title)
	_ = writer.WriteField("author", ab.Author)
	_ = writer.WriteField("series", ab.Series)
	_ = writer.WriteField("library", libraryID)
	_ = writer.WriteField("folder", folderID)

	// Add files to the request with progress reporting
	for i, file := range filesList {
		part, err := writer.CreateFormFile(strconv.Itoa(i), filepath.Base(file.Name()))
		if err != nil {
			return err
		}

		// Create a progress reader to track callback
		fileStat, err := file.Stat()
		if err != nil {
			return err
		}
		pr := &ProgressReader{
			FileId:   i,
			FileName: filepath.Base(file.Name()),
			Reader:   file,
			Size:     fileStat.Size(),
			Callback: callback,
		}

		_, err = io.Copy(part, pr)
		if err != nil {
			return err
		}
	}

	err := writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.url+"/api/upload", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.loginResponse.User.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload audiobook: %s", resp.Status)
	}

	return nil
}

// Progress Reader for file upload progress
type Fn func(fileId int, fileName string, size int64, pos int64, percent int)
type ProgressReader struct {
	FileId   int
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
