package slack

import (
	"github.com/nlopes/slack"
	"time"
)

const (
	// action is used for slack action and resolve response.
	actionSelectMonth = "selectMonth"
	actionNow         = "now"
	actionFast        = "fast"
	actionHelp        = "help"
	actionCustom      = "custom"
	actionCalendar    = "calendar"
)

var dt = time.Now()

var Calendar = slack.SectionBlock{
	Type: "section",
	Text: &slack.TextBlockObject{
		Type: "plain_text",
		Text: "Please,choose the date :calendar:",
	},
	BlockID:   actionCalendar,
	Accessory: Datepicker,
}

var Datepicker = slack.DatePickerBlockElement{
	Type:        "datepicker",
	ActionID:    actionCustom,
	InitialDate: dt.Format("2006-01-02"),
	Placeholder: &slack.TextBlockObject{
		Type: "plain_text",
		Text: "Select a date",
	},
}
var ListTimeAction = []slack.AttachmentAction{
	{
		Text:  "8:00",
		Value: "8:00",
	},
	{
		Text:  "9:00",
		Value: "9:00",
	},
	{
		Text:  "10:00",
		Value: "10:00",
	},
	{
		Text:  "11:00",
		Value: "11:00",
	},
	{
		Text:  "12:00",
		Value: "12:00",
	},
	{
		Text:  "13:00",
		Value: "13:00",
	},
	{
		Text:  "14:00",
		Value: "14:00",
	},
	{
		Text:  "15:00",
		Value: "15:00",
	},
	{
		Text:  "16:00",
		Value: "16:00",
	},
	{
		Text:  "17:00",
		Value: "17:00",
	},
	{
		Text:  "18:00",
		Value: "18:00",
	},
	{
		Text:  "19:00",
		Value: "19:00",
	},
	{
		Text:  "20:00",
		Value: "20:00",
	},
	{
		Text:  "8:30",
		Value: "8:30",
	},
	{
		Text:  "9:30",
		Value: "9:30",
	},
	{
		Text:  "10:30",
		Value: "10:30",
	},
	{
		Text:  "11:30",
		Value: "11:30",
	},
	{
		Text:  "12:30",
		Value: "12:30",
	},
	{
		Text:  "13:30",
		Value: "13:30",
	},
	{
		Text:  "14:30",
		Value: "14:30",
	},
	{
		Text:  "15:30",
		Value: "15:30",
	},
	{
		Text:  "16:30",
		Value: "16:30",
	},
	{
		Text:  "17:30",
		Value: "17:30",
	},
	{
		Text:  "18:30",
		Value: "18:30",
	},
	{
		Text:  "19:30",
		Value: "19:30",
	},
	{
		Text:  "20:30",
		Value: "20:30",
	},
}

var ListButtonAction = []slack.AttachmentAction{
	{
		Name:  actionNow,
		Text:  "Now",
		Type:  "button",
		Style: "primary",
	},
	{
		Name:  actionFast,
		Text:  "Fast",
		Type:  "button",
		Style: "primary",
	},
	{
		Name:  actionCustom,
		Text:  "Custom",
		Type:  "button",
		Style: "primary",
	},
	{
		Name:  actionHelp,
		Text:  "Help",
		Type:  "button",
		Style: "primary",
	},
}

var HelpList = []slack.AttachmentField{
	{
		Title: "Now",
		Value: "View available rooms now for 30 minutes (choose location).",
	},
	{
		Title: "Fast",
		Value: "View available rooms for today (choose time and duration).",
	},
	{
		Title: "Custom",
		Value: "View available rooms (choose date, time, duration, location, subject).",
	},
	{
		Title: "Help",
		Value: "Help - show help message.",
	},
}
