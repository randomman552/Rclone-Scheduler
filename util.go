package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/urfave/cli/v2"
)

// Function to run the given request and parse the JSON response into the given type
func DoRequest[T any](request *http.Request) *T {
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Println(err)
		return nil
	}

	var result *T

	buf := bytes.Buffer{}
	buf.ReadFrom(response.Body)

	// For successful responses
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		err = json.Unmarshal(buf.Bytes(), &result)

		if err != nil {
			log.Println(err)
			return nil
		}

		return result
	} else {
		logString := "Unsuccessful status code '" + strconv.FormatInt(int64(response.StatusCode), 10) + "' when querying '" + request.URL.String() + "'"

		log.Println(logString)
	}

	return nil
}

// Utility function to build an rclone client from the CLI context
func getRCloneClient(c *cli.Context) RCloneClient {
	protocol := c.String("rclone.protocol")
	host := c.String("rclone.host")
	port := c.String("rclone.port")

	rloneUrl := fmt.Sprintf("%s://%s:%s", protocol, host, port)

	// Create api clients
	return NewRCloneClient(rloneUrl)
}

// Utility function to get the backup cron schedule from the CLI context
func getBackupSchedule(c *cli.Context) string {
	return c.String("backup.schedule")
}

// Utility function to get the backup source path from the CLI context
func getBackupSourcePath(c *cli.Context) string {
	src := c.String("backup.source")

	return src
}

// Utility function to get the destination path from the CLI context
func getBackupDestinationPath(c *cli.Context) string {
	remote := c.String("backup.remote")
	path := c.String("backup.destination")

	return remote + "/" + path
}
