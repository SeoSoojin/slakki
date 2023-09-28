package slakki

import (
	"github.com/slack-go/slack"
)

type CMDHandler[C CMDTypes] func(client *slack.Client, command C) error

type CMDTypes interface {
	slack.SlashCommand | slack.InteractionCallback
}

type ErrorHandler func(client *slack.Client, channel string, err error) error
