package fwew_lib

import (
	"fmt"
	"strings"
)

// Errors raised by package x.
const (
	// cache
	DictionaryNotFound = constError("no dictionary found")
	// numbers
	NegativeNumber     = constError("negative numbers not allowed")
	NumberTooBig       = constError("number too big")
	NoTranslationFound = constError("no translation found")
	// list
	InvalidNumber = constError("invalidNumericError")
)

// errors are basically strings, that implement the error interface
type constError string

// Implement error interface, so this is a valid error.
func (err constError) Error() string {
	return string(err)
}

// Implement newer Is method to check if wrapped error is the desired error.
func (err constError) Is(target error) bool {
	targetError := target.Error()
	errorString := string(err)
	return targetError == errorString || strings.HasPrefix(targetError, errorString+": ")
}

// wrap suberror with this error. `Is` can be checked if wrapped errors is of type
func (err constError) wrap(inner error) error {
	return wrapError{msg: string(err), err: inner}
}

// If an error is wrapped, we change the type to this
type wrapError struct {
	err error
	msg string
}

// Also implement Error interface, to use wrapErrors as error
func (err wrapError) Error() string {
	if err.err != nil {
		return fmt.Sprintf("%s: %v", err.msg, err.err)
	}
	return err.msg
}

// Make it possible to Unwrap a wrapped error again.
func (err wrapError) Unwrap() error {
	return err.err
}

// Implement newer Is method to check unwrapped error
func (err wrapError) Is(target error) bool {
	return constError(err.msg).Is(target)
}
