package server

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dghubble/sling"
	"github.com/nlopes/slack"
	"golang.org/x/oauth2"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	markups "slack-bot/slack-bot/pkg/slack"
	"strings"
)

const (
	// action is used for slack attament action.
	actionSelectMonth = "selectMonth"
	actionStart       = "start"
	actionCancel      = "cancel"
	actionSelectDays  = "selectMonth"
	actionNow         = "now"
	actionFast        = "fast"
	actionHelp        = "help"
	actionCustom      = "custom"
	actionCalendar    = "calendar"
)

var message slack.InteractionCallback

type interactionHandler struct {
	slackClient       *slack.Client
	verificationToken string
}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println(jsonStr)
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
	case actionSelectMonth:
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
	case actionHelp:
		title := fmt.Sprintf("Help List")
		originalMessage := message.OriginalMessage
		/*if originalMessage.Text == "" {
			log.Println("text nil")
			break
		}*/
		originalMessage.Attachments[0].Text = title
		originalMessage.Attachments[0].Fields = markups.HelpList
		originalMessage.Attachments[0].Actions = markups.ListButtonAction

		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&originalMessage)
		//responseMessage(w, originalMessage, title, "")
		return
	case actionCalendar:
		originalMessage := message.OriginalMessage
		originalMessage.Attachments[0].Actions = markups.ListTimeAction
		_ = json.NewEncoder(w).Encode(&originalMessage)
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

type slackOauthRequestParams struct {
	ClientID     string `url:"client_id,omitempty"`
	ClientSecret string `url:"client_secret,omitempty"`
	Scope        string `url:"scope,omitempty"`
	RedirectUri  string `url:"redirect_uri,omitempty"`
}

func generateOAuthRequest() (request *http.Request, err error) {
	params := slackOauthRequestParams{
		ClientID:     os.Getenv("SLACK_CLIENT_ID"),
		ClientSecret: os.Getenv("SLACK_CLIENT_SECRET"),
		Scope:        "commands,users:read",
	}
	request, err = sling.New().Get("https://slack.com/oauth/authorize").QueryStruct(params).Request()
	fmt.Println(request)
	if err != nil {
		fmt.Println("Bad request", err)
	}
	return request, nil
}

func (s *Server) Auth(writer http.ResponseWriter, request *http.Request) {
	oAuthRequest, err := generateOAuthRequest()
	fmt.Println(oAuthRequest)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to create OAuth Token request: %v", err), 501)

		errorMessage := fmt.Sprintf("Failed to create OAuth Token request, parameters: %v", request)
		log.Fatal(errorMessage)
		http.Error(writer, errorMessage, 501)
		return
	}
	return
}

// addToSlack initializes the oauth process and redirects to Slack
func (s *Server) addToSlack(w http.ResponseWriter, r *http.Request) {
	_ = os.Setenv("SLACK_CLIENT_ID", "617863072727.609868118610")
	_ = os.Setenv("SLACK_CLIENT_SECRET", "a39e1d00bbe6ce9a88c191391108600c")
	// Just generate random state
	b := make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		writeError(w, 500, err.Error())
	}
	globalState := hex.EncodeToString(b)
	conf := &oauth2.Config{
		ClientID:     os.Getenv("SLACK_CLIENT_ID"),
		ClientSecret: os.Getenv("SLACK_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://slack.com/oauth/authorize",
			TokenURL: "https://slack.com/api/oauth.access", // not actually used here
		},
		RedirectURL: "http://localhost:9000/auth",
		Scopes:      []string{"client"},
	}
	url := conf.AuthCodeURL(globalState)
	fmt.Println(url)
	http.Redirect(w, r, url, http.StatusFound)
}
func writeError(w http.ResponseWriter, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(err))
}

// auth receives the callback from Slack, validates and displays the user information
func (s *Server) auth(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")
	errStr := r.FormValue("error")
	if errStr != "" {
		writeError(w, 401, errStr)
		return
	}
	if state == "" || code == "" {
		writeError(w, 400, "Missing state or code")
		return
	}
	client := http.Client{}
	token, _, err := slack.GetOAuthToken(&client, os.Getenv("SLACK_CLIENT_ID"), os.Getenv("SLACK_CLIENT_SECRET"), code, "")
	fmt.Println(token)
	if err != nil {
		fmt.Println("Don't set Bot User OAuth Access Token")
		return
	}
	s.chToken <- token
}

type data struct {
	header    string
	shortInfo string
}

// home displays the add-to-slack button
func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	text := data{
		header:    "Install SlackBot",
		shortInfo: "This Slack OAuth button allows users to install the RoomBot to their Slack workspace",
	}
	fmt.Println(text)
	tmpl, err := template.ParseFiles("slack-bot/pkg/server/SlackAuth.html")
	if err != nil {
		log.Fatal("Bad parse html")
	}
	err = tmpl.Execute(w, "")
}
