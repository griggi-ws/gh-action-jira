package jira

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/tomhjp/gh-action-jira/config"
)

func DoRequest(config config.JiraConfig, method, path string, query url.Values, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, generateURL(config.BaseURL, path, query), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.APIToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("API call %s %s failed (%d): %s", method, path, resp.StatusCode, string(bytes))
	}

	return bytes, nil
}

func generateURL(baseURL string, path string, query url.Values) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	url := fmt.Sprintf("%s%s", baseURL, path)
	queryString := query.Encode()
	if queryString != "" {
		url += fmt.Sprintf("?%s", queryString)
	}
	return url
}
