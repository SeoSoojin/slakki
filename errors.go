package slakki

import (
	"errors"
	"time"

	"github.com/slack-go/slack"
)

var (
	ErrInvalidCommand  = errors.New("invalid command")
	ErrNilClient       = errors.New("slack client is nil")
	ErrNilSocket       = errors.New("socket client is nil")
	ErrNilHandler      = errors.New("handler is nil")
	ErrNilManager      = errors.New("manager is nil")
	ErrInvalidManager  = errors.New("invalid manager")
	ErrCommandNotFound = errors.New("command not found")
)

func DefaultError(err error) slack.Attachment {

	return slack.Attachment{
		Fields: []slack.AttachmentField{
			{
				Title: "Date: ",
				Value: time.Now().Format("2006-01-02 15:04:05"),
			},
			{
				Title: "Message: ",
				Value: err.Error(),
			},
		},
		Color: "#EF1942",
		Title: "Error",
	}

}

func renderError(client *slack.Client, channel string, err error) error {
	_, _, err = client.PostMessage(channel, slack.MsgOptionAttachments(DefaultError(err)))
	return err
}
