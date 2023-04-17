package pubsub

import "strings"

type MultiError []error

func (me MultiError) Error() string {
	var errors []string
	for _, err := range me {
		errors = append(errors, err.Error())
	}

	return strings.Join(errors, "; ")
}
