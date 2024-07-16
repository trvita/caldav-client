package input

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	str, err := String(bytes.NewBufferString("test\n"), "enter: ")
	assert.Equal(t, "test", str)
	assert.NoError(t, err)

	str, err = String(bytes.NewBufferString("\n"), "enter:")
	assert.Empty(t, str)
	assert.NoError(t, err)
}

func TestStringFail(t *testing.T) {
	str, err := String(bytes.NewBufferString(""), "enter:")
	assert.Empty(t, str)
	assert.Error(t, err)
}

func TestInt(t *testing.T) {
	n, err := Int(bytes.NewBufferString("42\n"), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, 42, n)

	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	maxIntStr := fmt.Sprintf("%d\n", maxInt)
	minIntStr := fmt.Sprintf("%d\n", minInt)
	n, err = Int(bytes.NewBufferString(maxIntStr), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, maxInt, n)
	n, err = Int(bytes.NewBufferString(minIntStr), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, minInt, n)
}

func TestIntFail(t *testing.T) {
	n, err := Int(bytes.NewBufferString(""), "enter: ")
	assert.Error(t, err)
	assert.Zero(t, n)

	n, err = Int(bytes.NewBufferString("words\n"), "enter: ")
	assert.Error(t, err)
	assert.Zero(t, n)
}

func TestInts(t *testing.T) {
	nums, err := Ints(bytes.NewBufferString("1, 2, 3\n"), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(nums))
	assert.Equal(t, []int{1, 2, 3}, nums)

	nums, err = Ints(bytes.NewBufferString("1\n"), "enter: ")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(nums))
	assert.Equal(t, []int{1}, nums)
}

func TetIntsFail(t *testing.T) {
	nums, err := Ints(bytes.NewBufferString("1 2 3\n"), "enter: ")
	assert.Error(t, err)
	assert.Empty(t, nums)

	nums, err = Ints(bytes.NewBufferString(",\n"), "enter: ")
	assert.Error(t, err)
	assert.Empty(t, nums)

	nums, err = Ints(bytes.NewBufferString("1, 2, 3"), "enter: ")
	assert.Error(t, err)
	assert.Empty(t, nums)

	nums, err = Ints(bytes.NewBufferString("words\n"), "enter: ")
	assert.Error(t, err)
	assert.Empty(t, nums)
}

// testing Event is a problem since first call of String s all input and others just  eof
