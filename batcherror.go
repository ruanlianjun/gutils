package gutils

import (
	"bytes"
)

type (
	BatchError struct {
		errors errorArray
	}
	errorArray []error
)

func (b *BatchError) Add(err error) {
	b.errors = append(b.errors, err)
}

func (b *BatchError) Err() error {
	switch len(b.errors) {
	case 0:
		return nil
	case 1:
		return b.errors[0]
	default:
		return b.errors
	}
}

func (b *BatchError) NotNil() bool {
	return len(b.errors) > 0
}

func (a errorArray) Error() string {
	var buf bytes.Buffer

	for i, err := range a {
		if i > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(err.Error())
	}
	return buf.String()
}
