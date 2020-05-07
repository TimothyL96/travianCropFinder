package main

import (
	"strings"
)

type Error struct {
	error
	errorString string
}

func (e *Error) Set(errorStrings ...string) {
	var strBuilder strings.Builder

	for _, v := range errorStrings {
		strBuilder.WriteString(v)
		strBuilder.WriteString(" ")
	}
	completeString := strBuilder.String()
	completeString = completeString[:len(completeString)-1]

	e.errorString += completeString + "\n"
}

func (e Error) Error() string {
	err := e.errorString
	if e.error != nil {
		err += e.error.Error()
	}
	// log.Print(err)

	return err
}
