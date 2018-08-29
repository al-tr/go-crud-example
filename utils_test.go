package main

import (
	"errors"
	"testing"
)

func TestPanicerr(t *testing.T) {
	t.Run("panicerr thros the error if it's not nil", func(t *testing.T) {
		defer func() {
			errorText := recover().(error).Error()
			if errorText != "test passed" {
				t.Error("Something went wrong")
			}
		}()

		err := errors.New("test passed")
		panicerr(err)

		t.Error("Failure, should've thrown an error before")
	})
	t.Run("panicerr ignored nil error", func(t *testing.T) {
		defer func() {
			if recover() != nil {
				errRecovery := recover().(error)
				t.Error("Should not have an error ", errRecovery.Error())
			}
		}()
		panicerr(nil)
	})
}

