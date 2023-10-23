package slakki

type CommandOptions func(*commandConfig)

type commandConfig struct {
	name      string
	Callbacks map[string]CallbackHandler
	Help      HelpHandler
}

func WithCallback(id string, callback CallbackHandler) CommandOptions {
	return func(o *commandConfig) {
		o.Callbacks[id] = callback
	}
}

func WithHelp(help HelpHandler) CommandOptions {
	return func(o *commandConfig) {
		o.Help = help
	}
}

func commmandOptionsCompose(name string, opts ...CommandOptions) *commandConfig {
	config := &commandConfig{}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

func (c commandConfig) Apply(manager *manager) error {

	if c.Callbacks != nil {
		for id, callback := range c.Callbacks {
			manager.callback(id, callback)
		}
	}

	if c.Help != nil {
		manager.help(c.name, c.Help)
	}

	return nil

}
