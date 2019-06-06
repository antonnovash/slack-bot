package controller
/*
import (
	"fmt"
	"log"

	"github.com/Lighty0410/event-bot/pkg/exchange"
	"github.com/Lighty0410/event-bot/pkg/google/sheet"
	"github.com/Lighty0410/event-bot/pkg/google/youtube"
	"github.com/Lighty0410/event-bot/pkg/slack"
)

const sheetGeneral = "general"

// ToExchange creates a new meeting in Exchange calendar.
func (c *Controller) ToExchange(m *Meeting) error {
	meeting := &exchange.Meeting{
		Location: m.Location,
		Start:    m.Start,
		End:      m.End,
		Subject:  m.Subject,
		Body:     m.Description,
	}
	err := c.Exchange.AddMeeting(meeting)
	if err != nil {
		return fmt.Errorf("cannot add new event: %v", err)
	}
	return nil
}

// ToSheets creates a new row with meeting info in Google Sheet.
func (c *Controller) ToSheets(m *Meeting) error {
	meeting := &sheet.Meeting{
		Location:    m.Location,
		Start:       m.Start,
		End:         m.End,
		Subject:     m.Subject,
		Description: m.Description,
		URL:         m.YoutubeLink,
		Key:         m.YoutubeKey,
		ImageID:     m.ImageID,
		StreamID:    m.StreamID,
	}
	// adds meeting info in sheets with provided names.
	addMeeting := func(sheets ...string) error {
		for _, name := range sheets {
			sh, isNew, err := c.Drive.Sheet(name)
			if err != nil {
				return fmt.Errorf("cannot get sheet %s: %v", name, err)
			}
			if isNew {
				err = c.Sheets.AddColumnNames(sh.Id)
				if err != nil {
					return fmt.Errorf("cannot add columns name: %v", err)
				}
			}
			err = c.Sheets.AddRecord(sh.Id, meeting)
			if err != nil {
				return fmt.Errorf("cannot add record: %v", err)
			}
		}
		return nil
	}
	err := addMeeting(sheetGeneral, m.Kind)
	if err != nil {
		return fmt.Errorf("cannot add meetings to sheets: %v", err)
	}
	return nil
}

// ToSlack sends a message containing meeting info to Slack.
func (c *Controller) ToSlack(m *Meeting) error {
	meeting := &slack.Meeting{
		Channel:     m.Kind,
		Location:    m.Location,
		Start:       m.Start,
		End:         m.End,
		Subject:     m.Subject,
		Description: m.Description,
		YoutubeLink: m.YoutubeLink,
	}
	err := c.Slack.SendMeeting(meeting)
	if err != nil {
		return fmt.Errorf("cannot send meeting: %v", err)
	}
	return nil
}

// ToYoutube creates a new YouTube broadcast.
func (c *Controller) ToYoutube(m *Meeting) (*youtube.BroadcastInfo, error) {
	meeting := &youtube.Meeting{
		Subject:     m.Subject,
		Kind:        m.Kind,
		Start:       m.Start,
		End:         m.End,
		Description: m.Description,
	}
	sheet, isNew, err := c.Drive.Sheet(m.Kind)
	if err != nil {
		return nil, fmt.Errorf("cannot get sheet: %v", err)
	}
	if isNew {
		err = c.Sheets.AddColumnNames(sheet.Id)
		if err != nil {
			return nil, fmt.Errorf("cannot add columns name: %v", err)
		}
	}
	meeting.StreamID, err = c.Sheets.RetrieveStreamID(sheet.Id)
	if err != nil {
		log.Printf("cannot retrieve previous stream ID from sheets: %v", err)
	}
	if m.ImageID != "" {
		meeting.ImageURL = fmt.Sprintf("https://drive.google.com/uc?export=view&id=%s", m.ImageID)
	} else {
		id, err := c.Sheets.RetrieveImageID(sheet.Id)
		if err != nil {
			log.Printf("cannot retrieve image ID from sheets: %v", err)
		}
		meeting.ImageURL = fmt.Sprintf("https://drive.google.com/uc?export=view&id=%s", id)
	}
	bi, err := c.YouTube.CreateBroadcast(meeting)
	if err != nil {
		return nil, fmt.Errorf("cannot create broadcast: %v", err)
	}
	return bi, nil
}
*/