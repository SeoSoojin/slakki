package slakki

import (
	"context"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type Manager interface {
	Slash(name string, handler CMDHandler[slack.SlashCommand])
	Interaction(name string, handler CMDHandler[slack.InteractionCallback])
	Mount(prefix string, src Manager) (Manager, error)
	ListenAndServe() error
}

type manager struct {
	slashCommands      map[string]CMDHandler[slack.SlashCommand]
	interactionCommand map[string]CMDHandler[slack.InteractionCallback]
	sClient            *socketmode.Client
	client             *slack.Client
	errorHandler       ErrorHandler
}

func NewManager(sClient *socketmode.Client, client *slack.Client) Manager {
	return &manager{
		slashCommands:      make(map[string]CMDHandler[slack.SlashCommand]),
		interactionCommand: make(map[string]CMDHandler[slack.InteractionCallback]),
		sClient:            sClient,
		client:             client,
		errorHandler:       renderError,
	}
}

func (m *manager) SetErrorHandler(handler ErrorHandler) {
	m.errorHandler = handler
}

func (m *manager) Slash(name string, handler CMDHandler[slack.SlashCommand]) {
	m.slashCommands[name] = handler
}

func (m *manager) Interaction(name string, handler CMDHandler[slack.InteractionCallback]) {
	m.interactionCommand[name] = handler
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
			go m.handleEvent(event)
		}

	}(context.Background(), m.sClient)

	return m.sClient.Run()

}

func (m *manager) handleEvent(event socketmode.Event) {

	switch req := event.Data.(type) {

	case slack.SlashCommand:

		m.sClient.Ack(*event.Request)
		command := strings.TrimPrefix(req.Command, "/")
		cmd, ok := m.slashCommands[command]
		if !ok {
			m.errorHandler(m.client, req.ChannelID, ErrCommandNotFound)
		}

		err := cmd(m.client, req)
		if err != nil {
			m.errorHandler(m.client, req.ChannelID, err)
		}

	case slack.InteractionCallback:

		m.sClient.Ack(*event.Request)
		command := req.CallbackID
		cmd, ok := m.interactionCommand[command]
		if !ok {
			m.errorHandler(m.client, req.Channel.ID, ErrCommandNotFound)
		}

		err := cmd(m.client, req)
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

	for k, v := range srcM.interactionCommand {
		m.interactionCommand[prefix+k] = v
	}

	return m, nil

}