package core

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func (err *Error) Error() string {
	return err.Message
}

// New500FromError will build and return a `core.Error` instance from the given `error`.
//
// This routine will cast the error against any error types expected to come from
// this service's dependencies. It will always return an instance with a `500` status.
//
// NOTE: this function will also log any pertinent info related to the error.
func New500FromError(err error, log *logrus.Logger) *Error {
	err500 := NewError500()

	switch errType := err.(type) {
	case amqp.Error:
		log.WithFields(logrus.Fields{
			"type":        fmt.Sprintf("%T", errType),
			"reason":      errType.Reason,
			"code":        errType.Code,
			"recoverable": errType.Recover,
			"fromServer":  errType.Server,
		}).Error(errType.Error())
		return err500

	default:
		log.Errorf("Unknown error handled: %T: %s", errType, errType.Error())
		return err500
	}
}

// NewError500 will construct and return a vanilla 500 error.
func NewError500() *Error {
	return &Error{
		Message: "Internal server error. We're working on it.",
		Status:  500,
		Code:    "",
		Meta:    map[string]string{},
	}
}
