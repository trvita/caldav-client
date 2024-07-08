package ui

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	str, err := GetString(bytes.NewBufferString("test\n"), "enter: ")
	assert.Equal(t, str, "test")
	assert.NoError(t, err)

	str, err = GetString(bytes.NewBufferString("\n"), "enter:")
	assert.Equal(t, str, "")
	assert.NoError(t, err)
}

func TestGetStringFail(t *testing.T) {
	str, err := GetString(bytes.NewBufferString(""), "enter:")
	assert.Equal(t, str, "")
	assert.Error(t, err)
}

func TestGetInt(t *testing.T) {
	n, err := GetInt(bytes.NewBufferString("42\n"), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, n, 42)

	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	maxIntStr := fmt.Sprintf("%d\n", maxInt)
	minIntStr := fmt.Sprintf("%d\n", minInt)
	n, err = GetInt(bytes.NewBufferString(maxIntStr), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, n, maxInt)
	n, err = GetInt(bytes.NewBufferString(minIntStr), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, n, minInt)
}

func TestGetIntFail(t  *testing.T) {
	n, err := GetInt(bytes.NewBufferString(""), "enter: ")
	assert.Equal(t, n, 0)
	assert.Error(t, err)
	
	n, err = GetInt(bytes.NewBufferString("words\n"), "enter: ")
	assert.Equal(t, n, 0)
	assert.Error(t, err)
}

func TestGetInts(t *testing.T) {}

func TestGetUsernameBaikal(t *testing.T) {}

func TestGetEvent(t *testing.T) {}

func TestGetRecurrentEvent(t *testing.T) {

}
