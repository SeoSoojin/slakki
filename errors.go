package slakki

import "errors"

var (
	ErrInvalidCommand = errors.New("invalid command")
	ErrNilClient      = errors.New("slack client is nil")
	ErrNilSocket      = errors.New("socket client is nil")
	ErrNilHandler     = errors.New("handler is nil")
	ErrNilManager     = errors.New("manager is nil")
	ErrInvalidManager = errors.New("invalid manager")
)
