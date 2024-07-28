package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

// Struct represnenting an RClone job status request
type RCloneJobStatusRequest struct {
	JobId int `json:"jobid"`
}

// Struct representing an RClone job status response
type RCloneJobStatusReponse struct {
	Id        int       `json:"id"`
	Duration  float64   `json:"duration"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Error     string    `json:"error"`
	Finished  bool      `json:"finished"`
	Success   bool      `json:"success"`
}

// Struct representing an RClone core stats request
type RCloneCoreStatsRequest struct {
	Group string `json:"group"`
}

// Struct representing an RClone core stats response
type RCloneCoreStatsResponse struct {
	Bytes               uint    `json:"bytes"`
	Checks              uint    `json:"checks"`
	Deletes             uint    `json:"deletes"`
	Transfers           uint    `json:"transfers"`
	Errors              uint    `json:"errors"`
	Renames             uint    `json:"renames"`
	ElapsedTime         float64 `json:"elapsedTime"`
	Eta                 float64 `json:"eta"`
	FatalError          bool    `json:"fatalError"`
	LastError           string  `json:"lastError"`
	RetryError          bool    `json:"retryError"`
	ServerSideCopies    uint    `json:"serverSideCopies"`
	ServerSideCopyBytes uint    `json:"servierSideCopyBytes"`
	ServerSideMoves     uint    `json:"serverSideMoves"`
	ServerSideMoveBytes uint    `json:"serverSideMoveBytes"`
	Speed               float64 `json:"speed"`
	TransferTime        float64 `json:"transferTime"`
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
func (c RCloneClient) StartSync(source string, dest string) *RCloneSyncResponse {
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

// Get the status of the given sync job
func (c RCloneClient) GetSyncStatus(jobId int) *RCloneJobStatusReponse {
	bodyStruct := RCloneJobStatusRequest{
		JobId: jobId,
	}

	// Encode request
	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(bodyStruct)

	httpRequest := c.NewRequest("POST", "/job/status", body)
	response := DoRequest[RCloneJobStatusReponse](httpRequest)

	return response
}

// Get the transfer stats for the given sync job
func (c RCloneClient) GetSyncStats(jobId int) *RCloneCoreStatsResponse {
	bodyStruct := RCloneCoreStatsRequest{
		Group: fmt.Sprintf("job/%d", jobId),
	}

	// Encode request
	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(bodyStruct)

	httpRequest := c.NewRequest("POST", "/core/stats", body)
	response := DoRequest[RCloneCoreStatsResponse](httpRequest)

	return response
}
