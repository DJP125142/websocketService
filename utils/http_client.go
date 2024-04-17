package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

/*
* GET请求
 */
func SendGetRequest(baseURL string, params map[string]string, options ...string) (string, error) {
	// Prepare URL with parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	// Create a new request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	// Check if token was provided and add Authorization header
	if len(options) > 0 && options[0] != "" {
		req.Header.Add("Authorization", "Bearer "+options[0])
	}

	// Create a client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

/*
* POST请求
 */
func SendPostRequest(urlStr string, params map[string]interface{}) (string, error) {
	// Marshal the parameters to JSON for the body
	body, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	// Prepare the URL with query parameters
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for key, value := range params {
		// Convert each parameter value to string
		// Note: This conversion might need to be adjusted based on the expected format of query parameters
		strValue, ok := value.(string)
		if !ok {
			// Handle or skip non-string values
			continue // or return an error if necessary
		}
		q.Set(key, strValue)
	}
	u.RawQuery = q.Encode()

	// Create a new HTTP request with the JSON body
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseBody), nil
}
