package slack

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
)

const (
	// action is used for slack attament action.
	actionSelectMonth = "selectMonth"
	actionSelectDays = "selectMonth"
	actionStart  = "start"
	actionCancel = "cancel"
)

type SlackListener struct {
	client    *slack.Client
	channelID string
	rtm       *slack.RTM
}

func New(c Config) (*SlackListener, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("config is invalid: %v", err)
	}
	client := slack.New(c.Token)
	rtm := client.NewRTM()
	return &SlackListener{client: client, channelID: c.ChannelID, rtm: rtm}, nil
}
// handleMesageEvent handles message events.

// handleMesageEvent handles message events.
func (s *SlackListener) HandleMessageEvent(ev *slack.MessageEvent) error {

	// Only response in specific channel. Ignore else.
/*	if ev.Channel != s.channelID {
		log.Printf("%s %s", ev.Channel, ev.Msg.Text)
		return nil
	}*/
	// Only response mention to bot. Ignore else.
	/*if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.botID)) {
		return nil
	}*/
	botCommand := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	switch botCommand[0] {
	case "start":
		_ = s.HandleMessageStart(ev)


	/*case"custom":
		_ = s.HandleMessageCustom(ev)
*/
	}
	return nil
}
func (s *SlackListener) HandleMessageStart(ev *slack.MessageEvent) error {
	// value is passed to message handler when request is approved.
	attachmentMonth := slack.Attachment{
		Text:       "Chose the month,please:thinking_cat_face:",
		Color:      "#f9a41b",
		CallbackID: "Month",
		Actions: []slack.AttachmentAction{
			{
				Name: actionSelectMonth,
				Type: "select",
				Options: []slack.AttachmentActionOption{
					{
						Text:  "January",
						Value: "January",
					},
					{
						Text:  "February",
						Value: "February",
					},
					{
						Text:  "March",
						Value: "March",
					},
					{
						Text:  "April",
						Value: "April",
					},
					{
						Text:  "May",
						Value: "May",
					},
					{
						Text:  "June",
						Value: "June",
					},
					{
						Text:  "July",
						Value: "July",
					},
					{
						Text:  "Aughts",
						Value: "Aughts",
					},
					{
						Text:  "September",
						Value: "September",
					},
					{
						Text:  "October",
						Value: "October",
					},
					{
						Text:  "November",
						Value: "November",
					},
					{
						Text:  "December",
						Value: "December",
					},

				},

			},

			/*{
				Name:  actionStart,
				Text:  "start",
				Type:  "button",
				Style: "default",
			},
			{
				Name:  actionCancel,
				Text:  "Cancel",
				Type:  "button",
				Style: "danger",
			},*/
		},
	}
	attachmentDays := slack.Attachment{
		Text:       " Days:beer:",
		Color:      "#f9a41b",
		CallbackID: "beer",
		Actions: []slack.AttachmentAction{
			{
				Name: actionSelectDays,
				Type: "select",
				Options: []slack.AttachmentActionOption{
					{
						Text:  "Asahi Super Dry",
						Value: "Asahi Super Dry",
					},
					{
						Text:  "Kirin Lager Beer",
						Value: "Kirin Lager Beer",
					},
					{
						Text:  "Sapporo Black Label",
						Value: "Sapporo Black Label",
					},
					{
						Text:  "Suntory Malts",
						Value: "Suntory Malts",
					},
					{
						Text:  "Yona Yona Ale",
						Value: "Yona Yona Ale",
					},
				},

			},

			/*{
				Name:  actionStart,
				Text:  "start",
				Type:  "button",
				Style: "default",
			},
			{
				Name:  actionCancel,
				Text:  "Cancel",
				Type:  "button",
				Style: "danger",
			},*/
		},
	}

	paramsMonth := slack.MsgOptionAttachments(attachmentMonth)
	paramsDays := slack.MsgOptionAttachments(attachmentDays)
	if _, _, err := s.client.PostMessage(ev.Channel, slack.MsgOptionText("А что мы видим?", false), paramsMonth); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}
	if _, _, err := s.client.PostMessage(ev.Channel, slack.MsgOptionText("А что мы видим?", false), paramsDays); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}


	return nil
}

/*type interactionHandler struct {
	slackClient       *slack.Client
	verificationToken string
}

func (h interactionHandler) Handler(w http.ResponseWriter, r *http.Request) {
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
	if message.Token != h.verificationToken {
		log.Printf("[ERROR] Invalid token: %s", message.Token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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
*/