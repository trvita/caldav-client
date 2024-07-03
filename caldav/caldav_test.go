package caldav

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCredentials(t *testing.T) {
	input := bytes.NewBufferString("testuser\ntestpassword\n")
	username, password, err := GetCredentials(input)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", username)
	assert.Equal(t, "testpassword", password)
}

func TestCreateClient(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient("http://localhost:5232", bytes.NewBufferString("testuser\ntestpassword\n"))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
}
