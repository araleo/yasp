package git

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type GitConfig struct {
	Token     string
	ApiUrl    string
	ProjectId string
}

var gitConfig = GitConfig{}

// LoadConfig loads the environment variables into config variable.
func LoadConfig() {
	gitConfig.Token = os.Getenv("GITLAB_TOKEN")
	gitConfig.ApiUrl = os.Getenv("GITLAB_API_URL")
	gitConfig.ProjectId = os.Getenv("GITLAB_PROJECT_ID")
}

// sendRequest is an abstract function to send customizable requests with a token auth header.
func sendRequest(method string, url string, body io.Reader, token string) *http.Response {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("PRIVATE-TOKEN", token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}

// parseResp parses an http.Response and returns it's body as an slice of bytes.
func parseResp(resp *http.Response) []byte {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

// CreateIssue creates a new issue GitLab with the title param in the project specified in the env configuration and returns the new issue iid.
func CreateIssue(title string) int {
	endpoint := fmt.Sprintf("%s/%s/issues?title=%s", gitConfig.ApiUrl, gitConfig.ProjectId, title)
	resp := sendRequest("POST", endpoint, nil, gitConfig.Token)
	defer resp.Body.Close()
	body := parseResp(resp)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	return int(result["iid"].(float64))
}

// ListIssues lists the GitLab issues of the project specified in the env config.
func ListIssues() {
	endpoint := gitConfig.ApiUrl + "/" + gitConfig.ProjectId + "/issues"
	resp := sendRequest("GET", endpoint, nil, gitConfig.Token)
	defer resp.Body.Close()
	body := parseResp(resp)

	var result []map[string]interface{}
	json.Unmarshal(body, &result)

	for _, issue := range result {
		id := int(issue["id"].(float64))
		author := issue["author"].(map[string]interface{})

		fmt.Println()
		fmt.Println("Id: ", id)
		fmt.Println("Iid: ", issue["iid"])
		fmt.Println("Title: ", issue["title"])
		fmt.Println("Description: ", issue["description"])
		fmt.Println("Author: ", author["username"])
		fmt.Println("State: ", issue["state"])
	}

}
