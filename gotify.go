package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"text/template"

	"github.com/urfave/cli/v2"
)

// Struct to send notifications to gotify
type GotifyNotifier struct {
	// The token for gotify
	GotifyToken string
	// The url for gotify
	GotifyUrl string
}

// Struct containing extra display options for gotify
type MessageRequestExtrasDisplay struct {
	ContentType string `json:"contentType"`
}

// Struct containing extra options for gotify
type MessageRequestExtras struct {
	Display MessageRequestExtrasDisplay `json:"client::display"`
}

// Struct used to send a message request to gotify
type MessageRequest struct {
	Title    string               `json:"title"`
	Message  string               `json:"message"`
	Priority int                  `json:"priority"`
	Extras   MessageRequestExtras `json:"extras"`
}

// Create a new GotifyNotifer from the current cli context
func NewGotifyNotifier(c *cli.Context) GotifyNotifier {
	url := c.String("gotify.url")
	token := c.String("gotify.token")

	return GotifyNotifier{
		GotifyUrl:   url,
		GotifyToken: token,
	}
}

// Build the message URL from settings
func (g GotifyNotifier) GetMessageUrl() string {
	url, err := url.Parse(g.GotifyUrl)

	if err != nil {
		log.Printf("Error parsing gotify url: %s", err)
	}

	url = url.JoinPath("/message")

	// Setup query string
	query := url.Query()
	query.Set("token", g.GotifyToken)
	url.RawQuery = query.Encode()
	return url.String()
}

// Create a new gotify message request with the given body
func (g GotifyNotifier) NewRequest(title string, priority int, message *bytes.Buffer) *http.Request {
	url := g.GetMessageUrl()

	// Create request body
	messageString := message.String()
	requestBody := MessageRequest{
		Title:    title,
		Message:  messageString,
		Priority: 5,
		Extras: MessageRequestExtras{
			Display: MessageRequestExtrasDisplay{
				ContentType: "text/markdown",
			},
		},
	}
	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(requestBody)

	// Create request
	request, err := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Printf("Error creating request: %s", err)
		return nil
	}

	return request
}

func (g GotifyNotifier) RenderTemplate(path string, data any) *bytes.Buffer {
	tplate, err := template.ParseFiles(path)

	if err != nil {
		log.Printf("Error reading template: %s", err)
		return nil
	}

	result := &bytes.Buffer{}
	err = tplate.Execute(result, data)

	if err != nil {
		log.Printf("Error rendering template: %s", err)
		return nil
	}

	return result
}

// Check this notifyer is enabled
func (g GotifyNotifier) IsEnabled() bool {
	return len(g.GotifyUrl) > 0 && len(g.GotifyToken) > 0
}

// Execute the given request and handle any errors
func (g GotifyNotifier) DoRequest(request *http.Request) {
	client := http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Printf("Failed to notify gotify: %s", err)
		return
	}

	if response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		log.Printf("Gotify request error: %s - %s", response.Status, body)
	}

	response.Body.Close()
}

// Send a gotify notification that the backup has started
func (g GotifyNotifier) NotifyBackupStarted(context BackupStartedContext) {
	message := g.RenderTemplate("templates/started/message.tmpl", context)
	title := g.RenderTemplate("templates/started/title.tmpl", context)

	request := g.NewRequest(title.String(), 5, message)
	g.DoRequest(request)
	log.Printf("Sent gotify backup started notification")
}

func (g GotifyNotifier) NotifyBackupFinished(context BackupFinishedContext) {
	message := g.RenderTemplate("templates/finished/message.tmpl", context)
	title := g.RenderTemplate("templates/finished/title.tmpl", context)

	request := g.NewRequest(title.String(), 5, message)
	g.DoRequest(request)
	log.Printf("Sent gotify backup finished notification")
}
