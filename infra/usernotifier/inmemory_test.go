package usernotifier

import "testing"

func TestInMemory(t *testing.T) {
	runTestSuite(t, NewInMemory())
}
