package audiobookshelf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Call the Audiobookshelf API login method.
func Login(absUrl string, username, password string) (LoginResponse, error) {
	requestBody := LoginRequest{
		Username: username,
		Password: password,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("failed to marshal login response body: %v", err)
	}

	resp, err := http.Post(absUrl, "application/json", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return LoginResponse{}, fmt.Errorf("failed to make login API call: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LoginResponse{}, fmt.Errorf("login API call returned status code: %d", resp.StatusCode)
	}

	var loginResp LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("failed to decode login response: %v", err)
	}

	return loginResp, nil
}

// Call the Audiobookshelf API libraries method
func Libraries(url string, token string) (LibrariesResponse, error) {

	client := &http.Client{}
	url = url + "/api/libraries"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LibrariesResponse{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return LibrariesResponse{}, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return LibrariesResponse{}, fmt.Errorf("error reading response body: %v", err)
	}

	var response LibrariesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return LibrariesResponse{}, fmt.Errorf("error parsing response body: %v", err)
	}

	return response, nil
}

func GetLibraryByName(libraries []Library, libraryName string) (libraryID string, err error) {	

	for _, library := range libraries {
		if library.Name == libraryName {
			return library.ID, nil
		}
	}
	return "", fmt.Errorf("no library with name '%s' found", libraryName)
}

// Call the Audiobookshelf API
func ScanLibrary(url string, authToken string, libraryID string) error {
	url = url + "/api/libraries/" + libraryID + "/scan"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+authToken)

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