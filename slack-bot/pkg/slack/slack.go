package slack

import (
	"fmt"
	"github.com/nlopes/slack"
	"strings"
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
	if len(botCommand) == 0 {
		return nil
	}
	switch botCommand[0] {
	case "start":
		_ = s.HandleMessageStart(ev)
	case "custom":
		_ = s.HandleMessageCustom(ev)
	/*case "fast":
		_ = s.HandleMessageFast(ev)
	case "now":
		_ = s.HandleMessageNow(ev)*/
	case "help":
		_ = s.HandleMessageHelp(ev)
	default:
		return nil
	}
	return nil
}
func (s *SlackListener) HandleMessageCustom(ev *slack.MessageEvent) error {
	// value is passed to message handler when request is approved.
	/*attachmentMonth := slack.Attachment{
		Text:       "Month",
		Color:      "#72e004",
		CallbackID: "Month",
		Actions: []slack.AttachmentAction{
			{
				Name:    actionSelectMonth,
				Type:    "select",
				Options: ListMonthAction,
			},
		},
	}
	attachmentDays := slack.Attachment{
		Text:       " Days",
		Color:      "#72e004",
		CallbackID: "beer",
		Actions: []slack.AttachmentAction{
			{
				Name:    actionSelectDays,
				Type:    "select",
				Options: ListDaysAction,
			},
		},
	}

	paramsMonth := slack.MsgOptionAttachments(attachmentMonth)
	fmt.Println(paramsMonth)
	paramsDays := slack.MsgOptionAttachments(attachmentDays)
	if _, _, err := s.client.PostMessage(ev.Channel, paramsMonth); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}
	if _, _, err := s.client.PostMessage(ev.Channel, slack.MsgOptionText("", false), paramsDays); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}*/
	paramsCalendar := slack.MsgOptionBlocks(Calendar)
	fmt.Println(paramsCalendar)
	if _, _, err := s.client.PostMessage(ev.Channel, paramsCalendar); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}
	attachmentListOfButtons := slack.Attachment{
		CallbackID: "BlaBla",
		Actions:    ListButtonAction,
	}
	paramsCalendarButton := slack.MsgOptionAttachments(attachmentListOfButtons)
	fmt.Println(paramsCalendar)
	if _, _, err := s.client.PostMessage(ev.Channel, paramsCalendarButton); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}
	return nil
}
func (s *SlackListener) HandleMessageHelp(ev *slack.MessageEvent) error {
	// value is passed to message handler when request is approved.
	attachmentHelp := slack.Attachment{
		Text:       "Help list:",
		Color:      "#72e004",
		CallbackID: "help",
		Fields:     HelpList,
		Actions:    ListButtonAction,
	}

	paramsHelp := slack.MsgOptionAttachments(attachmentHelp)

	if _, _, err := s.client.PostMessage(ev.Channel, slack.MsgOptionText("", false), paramsHelp); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	return nil
}

func (s *SlackListener) HandleMessageStart(ev *slack.MessageEvent) error {
	// value is passed to message handler when request is approved.
	attachmentStart := slack.Attachment{
		Text:       "Hello,I'm ItechArt Bot:sunglasses:I'm represent the list of available rooms.",
		Color:      "#e00404",
		CallbackID: "start",
		Actions:    ListButtonAction,
	}

	paramsStart := slack.MsgOptionAttachments(attachmentStart)

	if _, _, err := s.client.PostMessage(ev.Channel, slack.MsgOptionText("", false), paramsStart); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	return nil
}
