package caldav

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCredentials(t *testing.T) {
	input := bytes.NewBufferString("testuser\ntestpassword\n")
	username, password := GetCredentials(input)
	assert.Equal(t, "testuser", username)
	assert.Equal(t, "testpassword", password)
}
