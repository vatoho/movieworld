package errorreview

import "errors"

var (
	ErrorNoReview = errors.New("user has not got review with such id")
	ErrorNoLogger = errors.New("there is no logger in context")
)
