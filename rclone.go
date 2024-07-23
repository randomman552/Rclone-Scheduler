package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// API Client for RClone
type RCloneClient struct {
	Url string
}

// Struct representing an RClone sync request
type RCloneSyncRequest struct {
	Async       bool   `json:"_async"`
	Source      string `json:"srcFs"`
	Destination string `json:"dstFs"`
}

// Struct representing an RClone sync response
type RCloneSyncResponse struct {
	JobId int `json:"jobid"`
}

// Create a new RClone API client pointing at the given url
func NewRCloneClient(url string) RCloneClient {
	return RCloneClient{
		Url: url,
	}
}

// Create a new request with the given parameters
func (c *RCloneClient) NewRequest(method string, url string, body io.Reader) *http.Request {
	url = c.Url + url
	request, err := http.NewRequest(method, url, body)

	if err != nil {
		panic(err)
	}

	request.Header.Set("Content-Type", "application/json")

	return request
}

// Start a sync operation from the given source to the given destination
func (c *RCloneClient) StartSync(source string, dest string) *RCloneSyncResponse {
	// Create request body
	bodyStruct := RCloneSyncRequest{
		Async:       true,
		Source:      source,
		Destination: dest,
	}

	// Encode request
	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(bodyStruct)

	httpRequest := c.NewRequest("POST", "/sync/sync", body)
	response := DoRequest[RCloneSyncResponse](httpRequest)

	return response
}
