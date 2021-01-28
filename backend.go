package go_github_s3

import (
	"github.com/vivekchandrabs/go_github_s3/utils"
	"strings"
)

type Session struct {
	RepositoryName string
	AccessToken    string
	UserName       string
}

type PutObjectInput struct {
	FileName    string
	FileContent []byte
	BucketName  string
	BranchName  string
}

type DeleteObjectInput struct {
	FilePath   string
	BranchName string
}

func (s *Session) PutObject(input *PutObjectInput) (string, error) {
	config := make(map[string]string)
	config["github_username"] = s.UserName
	config["github_repo_name"] = s.RepositoryName
	config["bucket_name"] = input.BucketName
	config["file_name"] = input.FileName

	uploadUrl := utils.EndpointBuilder(config)
	authorizationHeaders := utils.GetHeaders(s.AccessToken)

	filename := config["file_name"]

	for true {
		response := utils.Exists(uploadUrl, authorizationHeaders)
		if response {
			uploadUrl = utils.EndpointBuilder(config)
			newName := utils.GetName(filename)
			config["file_name"] = newName
			continue
		}

		break
	}

	payload := make(map[string]interface{})
	payload["message"] = input.FileName
	payload["committer"] = map[string]interface{}{
		"name":  "vivekchandrabs",
		"email": "vivek.chandra.301096@gmail.com",
	}
	payload["content"] = utils.Base64Converter(input.FileContent)
	payload["branch"] = input.BranchName

	response, err := utils.ExecRequest("PUT", authorizationHeaders, uploadUrl, payload)
	if err != nil {
		return "", err
	}

	content := response.Body["content"].(map[string]interface{})

	return content["download_url"].(string), nil
}

func (s *Session) DeleteObject(input *DeleteObjectInput) (string, error) {
	delimiter := "/" + input.BranchName + "/"
	filePathArray := strings.Split(input.FilePath, delimiter)
	imagePath := filePathArray[len(filePathArray)-1]

	config := make(map[string]string)
	config["github_username"] = s.UserName
	config["github_repo_name"] = s.RepositoryName
	config["bucket_name"] = ""
	config["file_name"] = imagePath

	url := utils.EndpointBuilder(config)
	authorizationHeaders := utils.GetHeaders(s.AccessToken)

	response, err := utils.ExecRequest("GET", authorizationHeaders, url, nil)
	if err != nil {
		return "", err
	}

	sha := response.Body["sha"].(string)

	payload := make(map[string]interface{})
	payload["message"] = response.Body["name"].(string)
	payload["committer"] = map[string]interface{}{
		"name":  "vivekchandrabs",
		"email": "vivek.chandra.301096@gmail.com",
	}
	payload["branch"] = input.BranchName
	payload["sha"] = sha

	response, err = utils.ExecRequest("DELETE", authorizationHeaders, url, payload)
	if err != nil {
		return "", err
	}

	return "Image deleted successfully", nil
}
