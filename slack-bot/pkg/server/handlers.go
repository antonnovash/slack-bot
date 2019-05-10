package server

import (
	"encoding/json"
	"fmt"
	"github.com/dghubble/sling"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	// action is used for slack attament action.
	actionSelect = "selectMonth"
	actionStart  = "start"
	actionCancel = "cancel"
)

type interactionHandler struct {
	slackClient       *slack.Client
	verificationToken string
}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("lel")
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		log.Printf("[ERROR] Failed to unespace request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var message slack.InteractionCallback
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Only accept message from slack with valid token
	/*if message.Token != h.verificationToken {
		log.Printf("[ERROR] Invalid token: %s", message.Token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}*/

	action := message.Actions[0]
	switch action.Name {
	case actionSelect:
		value := action.SelectedOptions[0].Value

		// Overwrite original drop down message.
		originalMessage := message.OriginalMessage
		originalMessage.Attachments[0].Text = fmt.Sprintf("OK to order %s ?", strings.Title(value))
		originalMessage.Attachments[0].Actions = []slack.AttachmentAction{
			{
				Name:  actionStart,
				Text:  "Yes",
				Type:  "button",
				Value: "start",
				Style: "primary",
			},
			{
				Name:  actionCancel,
				Text:  "No",
				Type:  "button",
				Style: "danger",
			},
		}

		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&originalMessage)
		return
	case actionStart:
		title := ":ok: your order was submitted! yay!"
		responseMessage(w, message.OriginalMessage, title, "")
		return
	case actionCancel:
		title := fmt.Sprintf(":x: @%s canceled the request", message.User.Name)
		responseMessage(w, message.OriginalMessage, title, "")
		return
	default:
		log.Printf("[ERROR] ]Invalid action was submitted: %s", action.Name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// responseMessage response to the original slackbutton enabled message.
// It removes button and replace it with message which indicate how bot will work
func responseMessage(w http.ResponseWriter, original slack.Message, title, value string) {
	original.Attachments[0].Actions = []slack.AttachmentAction{} // empty buttons
	original.Attachments[0].Fields = []slack.AttachmentField{
		{
			Title: title,
			Value: value,
			Short: false,
		},
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&original)
}

type team struct {
	ID    string `gorm:"primary_key"`
	Name  string
	Scope string `json:"scope"`
	Token string
}

type slackOauthRequestParams struct {
	ClientID     string `url:"client_id,omitempty"`
	ClientSecret string `url:"client_secret,omitempty"`
	Code         string `url:"code,omitempty"`
}
type slackOauthRequestResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TeamName    string `json:"team_name"`
	TeamID      string `json:"team_id"`
}

func generateOAuthRequest(code string) (request *http.Request, err error) {
	os.Setenv("SLACK_CLIENT_ID", "617863072727.609868118610")
	os.Setenv("SLACK_CLIENT_SECRET","a39e1d00bbe6ce9a88c191391108600c")
	params := slackOauthRequestParams{
		ClientID:     os.Getenv("SLACK_CLIENT_ID"),
		ClientSecret: os.Getenv("SLACK_CLIENT_SECRET")}

	request, err = sling.New().Get("https://slack.com/oauth/authorize").QueryStruct(params).Request()
	fmt.Println(request)
	return
}

func (s *Server) Auth(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()

		code := request.URL.Query().Get("code")
		oAuthRequest, err := generateOAuthRequest(code)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Failed to create OAuth Token request: %v", err), 501)

			errorMessage := fmt.Sprintf("Failed to create OAuth Token request, parameters: %v", request)
			log.Fatal(errorMessage)
			http.Error(writer, errorMessage, 501)
			return
		}
		var client = &http.Client{
			Timeout: time.Second * 10,
		}
		_, err = client.Do(oAuthRequest)
		return
}
