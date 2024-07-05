package caldav

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var URL = "http://127.0.0.1:90/dav.php"
var testCredentials = "testuser\ntestpassword\n"

var listCalendarsOutput = "Calendar: cal-empty\nCalendar: cal-with-recs\nCalendar: cal-with-todos\nCalendar: default\n"
var listCalendarsOutputWithNew = "Calendar: cal-empty\nCalendar: cal-new\nCalendar: cal-with-recs\nCalendar: cal-with-todos\nCalendar: default\n"

var newCalendarName = "cal-new"
var existingCalendarName = "default"
var emptyCalendarName = "cal-empty"
var nonExistingCalendarName = "wrong"

var validUID = "uid"
var invalidUID = "invalid"

func TestGetCredentials(t *testing.T) {
	input := bytes.NewBufferString("testuser\ntestpassword\n")
	username, password, err := GetCredentials(input)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", username)
	assert.Equal(t, "testpassword", password)
}

func TestCreateClient(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
}

func TestExtractNameFromEmail(t *testing.T) {
	assert.Equal(t, "mr", ExtractNameFromEmail("mr@mail.com"))
	assert.Equal(t, "go", ExtractNameFromEmail("go@google"))
	assert.Equal(t, "", ExtractNameFromEmail("@"))
	assert.Equal(t, "", ExtractNameFromEmail("@mail.com"))
	assert.Equal(t, "mr", ExtractNameFromEmail("mr@"))
	assert.Equal(t, "", ExtractNameFromEmail("no-at"))
}

func TestListCalendars(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)
	err = ListCalendars(ctx, client, homeset)
	assert.NoError(t, err)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.Equal(t, output, listCalendarsOutput)
}

func TestCreateCalendar(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = CreateCalendar(ctx, httpClient, URL, homeset, newCalendarName, "0")
	assert.NoError(t, err)

	err = ListCalendars(ctx, client, homeset)
	assert.NoError(t, err)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.Equal(t, output, listCalendarsOutputWithNew)
	err = Delete(ctx, client, homeset+newCalendarName)
	assert.NoError(t, err)

}

func TestFindCalendarCorrect(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)

	err = FindCalendar(ctx, client, homeset, existingCalendarName)
	assert.NoError(t, err)
}

func TestFindCalendarWrong(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)

	err = FindCalendar(ctx, client, homeset, nonExistingCalendarName)
	assert.Error(t, err)
}

func TestListEvents(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)

	err = ListEvents(ctx, client, homeset, existingCalendarName)
	assert.NoError(t, err)
}

func TestListEventsWrong(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)

	err = ListEvents(ctx, client, homeset, emptyCalendarName)
	assert.Error(t, err)
}

func TestListTodos(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)

	err = ListEvents(ctx, client, homeset, existingCalendarName)
	assert.NoError(t, err)
}

func TestListTodosWrong(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)

	err = ListEvents(ctx, client, homeset, emptyCalendarName)
	assert.Error(t, err)
}

func TestCreateEvent(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)
	e := &Event{
		Name:          "VEVENT",
		Summary:       "event",
		Uid:           validUID,
		DateTimeStart: time.Now(),
		DateTimeEnd:   time.Now(),
		Attendees:     nil,
		Organizer:     "",
	}
	event := GetEvent(e)
	err = CreateEvent(ctx, client, homeset, existingCalendarName, event)
	assert.NoError(t, err)
}

func TestCreateEventWrong(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)
	e := &Event{}
	event := GetEvent(e)
	err = CreateEvent(ctx, client, homeset, existingCalendarName, event)
	assert.Error(t, err)
}
func TestDeleteEventWrong(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)
	err = Delete(ctx, client, homeset+existingCalendarName+"/"+invalidUID+".ics")
	assert.Error(t, err)
}
func TestDeleteCalendarWrong(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)
	err = Delete(ctx, client, homeset+nonExistingCalendarName)
	assert.Error(t, err)
}

// func TestFindEventsWithExpand(t *testing.T) {
// 	start, err := time.Parse("2006.01.02 15.04.05", "2020.02.02 00.00.00")
// 	assert.NoError(t, err)
// 	end, err := time.Parse("2006.01.02 15.04.05", "2025.02.02 00.00.00")
// 	assert.NoError(t, err)
// 	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials))
// 	assert.NoError(t, err)
// 	homeset, err := client.FindCalendarHomeSet(ctx, principal)
// 	assert.NoError(t, err)
// 	err = FindEventsWithExpand(ctx, httpClient, URL, homeset, "walks", start, end)
// 	assert.NoError(t, err)
// }
