package userstore

import "testing"

func TestInMemory(t *testing.T) {
	runTestSuite(t, NewInMemory())
}
