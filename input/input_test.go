package input

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputString(t *testing.T) {
	str, err := InputString(bytes.NewBufferString("test\n"), "enter: ")
	assert.Equal(t, "test", str)
	assert.NoError(t, err)

	str, err = InputString(bytes.NewBufferString("\n"), "enter:")
	assert.Empty(t, str)
	assert.NoError(t, err)
}

func TestInputStringFail(t *testing.T) {
	str, err := InputString(bytes.NewBufferString(""), "enter:")
	assert.Empty(t, str)
	assert.Error(t, err)
}

func TestInputInt(t *testing.T) {
	n, err := InputInt(bytes.NewBufferString("42\n"), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, 42, n)

	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	maxIntStr := fmt.Sprintf("%d\n", maxInt)
	minIntStr := fmt.Sprintf("%d\n", minInt)
	n, err = InputInt(bytes.NewBufferString(maxIntStr), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, maxInt, n)
	n, err = InputInt(bytes.NewBufferString(minIntStr), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, minInt, n)
}

func TestInputIntFail(t *testing.T) {
	n, err := InputInt(bytes.NewBufferString(""), "enter: ")
	assert.Error(t, err)
	assert.Zero(t, n)

	n, err = InputInt(bytes.NewBufferString("words\n"), "enter: ")
	assert.Error(t, err)
	assert.Zero(t, n)
}

func TestInputInts(t *testing.T) {
	nums, err := InputInts(bytes.NewBufferString("1, 2, 3\n"), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(nums))
	assert.Equal(t, []int{1, 2, 3}, nums)

	nums, err = InputInts(bytes.NewBufferString("1\n"), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(nums))
	assert.Equal(t, []int{1}, nums)
}

func TetInputIntsFail(t *testing.T) {
	nums, err := InputInts(bytes.NewBufferString("1 2 3\n"), "enter: ")
	assert.Error(t, err)
	assert.Empty(t, nums)

	nums, err = InputInts(bytes.NewBufferString(",\n"), "enter: ")
	assert.Error(t, err)
	assert.Empty(t, nums)

	nums, err = InputInts(bytes.NewBufferString("1, 2, 3"), "enter: ")
	assert.Error(t, err)
	assert.Empty(t, nums)

	nums, err = InputInts(bytes.NewBufferString("words\n"), "enter: ")
	assert.Error(t, err)
	assert.Empty(t, nums)
}

// testing InputEvent is a problem since first call of InputString Inputs all input and others just Input eof
