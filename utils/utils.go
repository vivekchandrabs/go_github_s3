package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
)

type ResponseData struct {
	StatusCode int
	Header     http.Header
	Body       map[string]interface{}
}

func EndpointBuilder(config map[string]string) string {
	url := "https://api.github.com/repos/"
	url = url + config["github_username"] + "/"
	url = url + config["github_repo_name"] + "/"
	url = url + "contents/"

	if config["bucket_name"] != "" {
		url = url + config["bucket_name"] + "/"
	}

	url = url + config["file_name"]

	return url
}

func Base64Converter(input []byte) string {
	//b := make([]byte, base64.StdEncoding.EncodedLen(len(input)))
	return base64.StdEncoding.EncodeToString(input)
}

func GetHeaders(accessToken string) string {
	return fmt.Sprintf("token %s", accessToken)
}

func GetName(name string) string {
	return uuid.New().String() + "-" + name
}

func Exists(path, authorizationHeader string) bool {
	response, err := ExecRequest("GET", authorizationHeader, path, nil)
	if err != nil {
		return false
	}

	if response.StatusCode == 200 {
		return true
	}

	return false
}

func ExecRequest(method, authorizationHeaders, url string, payload map[string]interface{}) (*ResponseData, error) {
	jsonString, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonString))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authorizationHeaders)

	// Initialize http client
	httpClient := http.Client{}

	// Execute http request
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Create response data struct
	responseData := &ResponseData{
		StatusCode: response.StatusCode,
		Header:     response.Header,
		Body:       make(map[string]interface{}),
	}

	// Convert response body type to byte slice
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Convert body type to map
	err = json.Unmarshal(body, &responseData.Body)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}
