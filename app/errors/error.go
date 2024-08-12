package errorapp

import "errors"

var (
	ErrorUserExists  = errors.New("user with such username already exist")
	ErrorNoFilm      = errors.New("film with such id does not exist")
	ErrorNoSession   = errors.New("no session with such id")
	ErrorNoLogger    = errors.New("no logger in context")
	ErrorNoRequestID = errors.New("no request id in logger")
)
