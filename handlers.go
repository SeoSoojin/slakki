package slakki

import "github.com/slack-go/slack"

type CMDHandler[C CMDTypes] func(client *slack.Client, command C) error

type CMDTypes interface {
	slack.SlashCommand | slack.InteractionCallback
}
