package errorauth

import "errors"

var (
	ErrorUserNotExist      = errors.New("no user with such username")
	ErrorBadPassword       = errors.New("wrong password for this user")
	ErrorUserAlreadyExists = errors.New("user with such username already exist")
	ErrorNoLogger          = errors.New("there is no logger in context")
)
