package ui

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	str, err := GetString(bytes.NewBufferString("test\n"), "enter: ")
	assert.Equal(t, str, "test\n")
	assert.NoError(t, err)
}
func TestGetStrings(t *testing.T) {}

func TestGetInt(t *testing.T) {}

func TestGetInts(t *testing.T) {}

func TestGetUsernameBaikal(t *testing.T) {}

func TestGetEvent(t *testing.T) {}

func TestGetRecurrentEvent(t *testing.T) {

}
