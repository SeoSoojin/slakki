package slakki

import (
	"context"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type Manager interface {
	Slash(name string, handler CMDHandler[slack.SlashCommand], opts ...CommandOptions)
	Callback(name string, handler CMDHandler[slack.InteractionCallback])
	SetErrorHandler(handler ErrorHandler)
	Mount(prefix string, src Manager) (Manager, error)
	ListenAndServe() error
}

type manager struct {
	slashCommands    map[string]SlashHandler
	callbackCommands map[string]CallbackHandler
	helpers          map[string]HelpHandler
	sClient          *socketmode.Client
	client           *slack.Client
	errorHandler     ErrorHandler
}

func NewManager(sClient *socketmode.Client, client *slack.Client) Manager {
	return &manager{
		slashCommands:    make(map[string]SlashHandler),
		callbackCommands: make(map[string]CallbackHandler),
		sClient:          sClient,
		client:           client,
		errorHandler:     renderError,
	}
}

func (m *manager) SetErrorHandler(handler ErrorHandler) {
	m.errorHandler = handler
}

func (m *manager) Slash(name string, handler CMDHandler[slack.SlashCommand], opts ...CommandOptions) {
	config := commmandOptionsCompose(name, opts...)
	config.Apply(m)
	m.slashCommands[name] = SlashHandler(handler)
}

func (m *manager) Callback(name string, handler CMDHandler[slack.InteractionCallback]) {
	m.callback(name, CallbackHandler(handler))
}

func (m *manager) callback(name string, handler CallbackHandler) {
	m.callbackCommands[name] = handler
}

func (m *manager) help(name string, handler HelpHandler) {
	m.helpers[name] = handler
}

func (m *manager) ListenAndServe() error {

	if m.client == nil {
		return ErrNilClient
	}

	if m.sClient == nil {
		return ErrNilSocket
	}

	go func(ctx context.Context, client *socketmode.Client) {

		for event := range client.Events {
			go m.handleEvent(ctx, event)
		}

	}(context.Background(), m.sClient)

	return m.sClient.Run()

}

func (m *manager) handleEvent(ctx context.Context, event socketmode.Event) {

	switch req := event.Data.(type) {

	case slack.SlashCommand:

		m.sClient.Ack(*event.Request)
		command := strings.TrimPrefix(req.Command, "/")
		if req.Text == "--help" {
			if helper, ok := m.helpers[command]; ok {
				helper(m.client, req.ChannelID, command)
				return
			}
		}

		cmd, ok := m.slashCommands[command]
		if !ok {
			m.errorHandler(m.client, req.ChannelID, ErrCommandNotFound)
		}

		err := cmd(ctx, m.client, req)
		if err != nil {
			m.errorHandler(m.client, req.ChannelID, err)
		}

	case slack.InteractionCallback:

		m.sClient.Ack(*event.Request)
		command := req.CallbackID
		cmd, ok := m.callbackCommands[command]
		if !ok {
			m.errorHandler(m.client, req.Channel.ID, ErrCommandNotFound)
		}

		err := cmd(ctx, m.client, req)
		if err != nil {
			m.errorHandler(m.client, req.Channel.ID, err)
		}

	}

}

func (m *manager) Mount(prefix string, src Manager) (Manager, error) {

	if src == nil {
		return nil, ErrNilManager
	}

	srcM, ok := src.(*manager)
	if !ok {
		return nil, ErrInvalidManager
	}

	prefix = strings.Trim(prefix, " ")

	for k, v := range srcM.slashCommands {
		m.slashCommands[prefix+k] = v
	}

	for k, v := range srcM.callbackCommands {
		m.callbackCommands[prefix+k] = v
	}

	return m, nil

}
