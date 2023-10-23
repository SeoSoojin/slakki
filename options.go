package slakki

type CommandOptions func(*commandConfig)

type commandConfig struct {
	name      string
	callbacks map[string]CallbackHandler
	help      HelpHandler
}

func WithCallback(id string, callback CallbackHandler) CommandOptions {
	return func(o *commandConfig) {
		o.callbacks[id] = callback
	}
}

func WithHelp(help HelpHandler) CommandOptions {
	return func(o *commandConfig) {
		o.help = help
	}
}

func commmandOptionsCompose(name string, opts ...CommandOptions) *commandConfig {
	config := &commandConfig{
		name:      name,
		callbacks: make(map[string]CallbackHandler),
	}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

func (c commandConfig) Apply(manager *manager) error {

	if c.callbacks != nil {
		for id, callback := range c.callbacks {
			manager.callback(id, callback)
		}
	}

	if c.help != nil {
		manager.help(c.name, c.help)
	}

	return nil

}
