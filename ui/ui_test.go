package ui

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	str, err := GetString(bytes.NewBufferString("test\n"), "enter: ")
	assert.Equal(t, "test", str)
	assert.NoError(t, err)

	str, err = GetString(bytes.NewBufferString("\n"), "enter:")
	assert.Empty(t, str)
	assert.NoError(t, err)
}

func TestGetStringFail(t *testing.T) {
	str, err := GetString(bytes.NewBufferString(""), "enter:")
	assert.Empty(t, str)
	assert.Error(t, err)
}

func TestGetInt(t *testing.T) {
	n, err := GetInt(bytes.NewBufferString("42\n"), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, 42, n)

	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	maxIntStr := fmt.Sprintf("%d\n", maxInt)
	minIntStr := fmt.Sprintf("%d\n", minInt)
	n, err = GetInt(bytes.NewBufferString(maxIntStr), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, maxInt, n)
	n, err = GetInt(bytes.NewBufferString(minIntStr), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, minInt, n)
}

func TestGetIntFail(t *testing.T) {
	n, err := GetInt(bytes.NewBufferString(""), "enter: ")
	assert.Error(t, err)
	assert.Zero(t, n)

	n, err = GetInt(bytes.NewBufferString("words\n"), "enter: ")
	assert.Error(t, err)
	assert.Zero(t, n)
}

func TestGetInts(t *testing.T) {
	nums, err := GetInts(bytes.NewBufferString("1, 2, 3\n"), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(nums))
	assert.Equal(t, []int{1, 2, 3}, nums)

	nums, err = GetInts(bytes.NewBufferString("1\n"), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(nums))
	assert.Equal(t, []int{1}, nums)
}

func TetGetIntsFail(t *testing.T) {
	nums, err := GetInts(bytes.NewBufferString("1 2 3\n"), "enter: ")
	assert.Error(t, err)
	assert.Empty(t, nums)

	nums, err = GetInts(bytes.NewBufferString(",\n"), "enter: ")
	assert.Error(t, err)
	assert.Empty(t, nums)

	nums, err = GetInts(bytes.NewBufferString("1, 2, 3"), "enter: ")
	assert.Error(t, err)
	assert.Empty(t, nums)

	nums, err = GetInts(bytes.NewBufferString("words\n"), "enter: ")
	assert.Error(t, err)
	assert.Empty(t, nums)
}

func TestGetUsernameBaikal(t *testing.T) {
	expectedName := "baikal"
	str := fmt.Sprintf("%s%s/", URLstart, expectedName)
	username := GetUsernameBaikal(str)
	assert.Equal(t, expectedName, username)

	str = fmt.Sprintf("%s/", expectedName)
	username = GetUsernameBaikal(str)
	assert.Equal(t, expectedName, username)

	username = GetUsernameBaikal("")
	assert.Empty(t, username)
}

func TestGetEvent(t *testing.T) {}

func TestGetRecurrentEvent(t *testing.T) {

}
