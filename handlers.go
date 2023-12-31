package slakki

import (
	"context"

	"github.com/slack-go/slack"
)

type SlashHandler func(ctx context.Context, client *slack.Client, command slack.SlashCommand) error
type CallbackHandler func(ctx context.Context, client *slack.Client, command slack.InteractionCallback) error

type ErrorHandler func(client *slack.Client, channel string, err error) error
type HelpHandler func(client *slack.Client, channel string, command string) error

type CMDTypes interface {
	slack.SlashCommand | slack.InteractionCallback
}

type CMDHandler[C CMDTypes] func(ctx context.Context, client *slack.Client, command C) error
