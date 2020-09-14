package errors

import (
	"fmt"
	"strings"
)

var (
	Unexpected = NewErr("Unexpected")

	WordDuplicate = NewErr("DuplicateWord")
	WordNotFound  = NewErr("WordNotFound")

	InvalidUUID         = NewErr("InvalidUUID")
	InvalidSessionToken = NewErr("InvalidSessionToken")
	InvalidSlug         = NewErr("InvalidSlug")
)

// -----------------------------------------------------------------------------
type errFace interface {
	Code() string
}

type Error struct {
	code string
	msg  string
}

func NewErr(code string) Error {
	return Error{code: code}
}

func (e Error) Error() string {
	if e.msg != "" {
		return fmt.Sprintf(strings.Join([]string{e.code, e.msg}, " - "))
	} else {
		return e.code
	}
}

func (e Error) Code() string {
	return e.code
}

func (e Error) Msg() string {
	return e.msg
}

func (e Error) WithMsg(msg string) error {
	return Error{e.code, msg}
}

func (e Error) WithErr(err error) error {
	return Error{e.code, err.Error()}
}

func (e Error) Equals(err error) bool {
	if err == nil {
		return false
	}
	e2, ok := err.(errFace)
	if !ok {
		return false
	}
	return e.code == e2.Code()
}

func UnexpectedError(err error, format string, args ...interface{}) error {
	//defer alerts.AlertError(err, format, args...)
	return Unexpected.WithErr(err)
}
