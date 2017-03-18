package testutils

import "testing"

type MockLog struct {
	Test *testing.T
}

func (log MockLog) Info(data string) string {
	return data
}

func (log MockLog) Error(data string) string {
	return data
}
