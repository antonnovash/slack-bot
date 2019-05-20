package slack

import (
	"github.com/nlopes/slack"
	"time"
)

const (
	// action is used for slack attament action.
	actionSelectMonth = "selectMonth"
	actionSelectDays  = "selectMonth"
	actionNow         = "now"
	actionFast        = "fast"
	actionHelp        = "help"
	actionCustom      = "custom"
	actionStart       = "start"
)

/*[
{
"type": "section",
"text": {
"type": "mrkdwn",
"text": "Pick a date for the deadline."
},
"accessory": {
"type": "datepicker",
"initial_date": "1990-04-28",
"placeholder": {
"type": "plain_text",
"text": "Select a date",
"emoji": true
}
}
}
]*/
var t = time.Now()
var block = slack.ActionBlock{

}
var datepicker = slack.DatePickerBlockElement{
	Type:        "datepicker",
	ActionID:    actionCustom,
	InitialDate: "1990-04-28",
	Placeholder: &slack.TextBlockObject{
		Type: "plain_text",
		Text: "Select a date",
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
var ListMonthAction = []slack.AttachmentActionOption{
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
}
var ListDaysAction = []slack.AttachmentActionOption{
	{
		Text:  "1",
		Value: "1",
	},
	{
		Text:  "2",
		Value: "2",
	},
	{
		Text:  "3",
		Value: "3",
	},
	{
		Text:  "4",
		Value: "4",
	},
	{
		Text:  "5",
		Value: "5",
	},
	{
		Text:  "6",
		Value: "6",
	},
	{
		Text:  "7",
		Value: "7",
	},
	{
		Text:  "8",
		Value: "8",
	},
	{
		Text:  "9",
		Value: "9",
	},
	{
		Text:  "10",
		Value: "10",
	},
	{
		Text:  "11",
		Value: "11",
	},
	{
		Text:  "12",
		Value: "12",
	},
	{
		Text:  "13",
		Value: "13",
	},
	{
		Text:  "14",
		Value: "14",
	},
	{
		Text:  "15",
		Value: "15",
	},
	{
		Text:  "16",
		Value: "16",
	},
	{
		Text:  "17",
		Value: "17",
	},
	{
		Text:  "18",
		Value: "18",
	},
	{
		Text:  "19",
		Value: "19",
	},
	{
		Text:  "20",
		Value: "20",
	},
	{
		Text:  "21",
		Value: "21",
	},
	{
		Text:  "22",
		Value: "22",
	},
	{
		Text:  "23",
		Value: "23",
	},
	{
		Text:  "24",
		Value: "24",
	},
	{
		Text:  "25",
		Value: "25",
	},
	{
		Text:  "26",
		Value: "26",
	},
	{
		Text:  "27",
		Value: "27",
	},
	{
		Text:  "28",
		Value: "28",
	},
	{
		Text:  "29",
		Value: "29",
	},
	{
		Text:  "30",
		Value: "30",
	},
	{
		Text:  "31",
		Value: "31",
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
