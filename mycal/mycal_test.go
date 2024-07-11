package mycal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/trvita/caldav-client-yandex/caldav"
	"github.com/trvita/go-ical"
)

var URL = "http://127.0.0.1:90/dav.php"
var testCredentials1 = "testuser\ntestpassword\n"
var testCredentials2 = "usertest\npasswordtest\n"

var testEmail1 = "some-mail@mail.com"
var testEmail2 = "mail-some@mail.com"

var listCalendarsOutput = "Calendar: cal-empty\nCalendar: cal-with-recs\nCalendar: cal-with-todos\nCalendar: default\n"
var listCalendarsOutputWithNew = "Calendar: cal-empty\nCalendar: cal-new\nCalendar: cal-with-recs\nCalendar: cal-with-todos\nCalendar: default\n"

var newCalendarName = "cal-new"
var existingCalendarName = "default"
var emptyCalendarName = "cal-empty"
var nonExistingCalendarName = "wrong"
var inboxCalendarName = "inbox"
var modCalendarName = "mod"

var validUID = "valid"
var invalidUID = "invalid"
var modificateUID = "modificate" // depends on user sending invitation

func TestGetCredentials(t *testing.T) {
	input := bytes.NewBufferString("testuser\ntestpassword\n")
	username, password, err := GetCredentials(input)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", username)
	assert.Equal(t, "testpassword", password)
}

func TestCreateClient(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials1))
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

	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials1))
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
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials1))
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
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials1))
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

func TestFindCalendarFail(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials1))
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

func TestGetEvents(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials1))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)

	resp, err := GetEvents(ctx, client, homeset, existingCalendarName)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}

func TestGetEventsFail(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials1))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)

	resp, err := GetEvents(ctx, client, homeset, emptyCalendarName)
	assert.Error(t, err)
	assert.Empty(t, resp)
}

func TestGetTodos(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials1))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)

	resp, err := GetEvents(ctx, client, homeset, existingCalendarName)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}

func TestGetTodosFail(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials1))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)

	resp, err := GetEvents(ctx, client, homeset, emptyCalendarName)
	assert.Error(t, err)
	assert.Empty(t, resp)
}

func TestCreateEvent(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials1))
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
	event, err := GetEvent(e)
	assert.NoError(t, err)
	err = CreateEvent(ctx, client, homeset, existingCalendarName, event)
	assert.NoError(t, err)
}

func TestCreateEventFail(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials1))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)
	e := &Event{}
	event, err := GetEvent(e)
	assert.NoError(t, err)
	err = CreateEvent(ctx, client, homeset, existingCalendarName, event)
	assert.Error(t, err)
}
func TestDeleteEventFail(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials1))
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
func TestDeleteCalendarFail(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials1))
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
// 	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(testCredentials1))
// 	assert.NoError(t, err)
// 	homeset, err := client.FindCalendarHomeSet(ctx, principal)
// 	assert.NoError(t, err)
// 	err = FindEventsWithExpand(ctx, httpClient, URL, homeset, "walks", start, end)
// 	assert.NoError(t, err)
// }

func TestAttend_CreateEventWithAttendee(t *testing.T) {
	var client *caldav.Client
	var principal, homeset string
	var ctx context.Context
	var err error
	// testuser creates event
	_, client, principal, ctx, err = CreateClient(URL, bytes.NewBufferString(testCredentials1))
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err = client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)
	assert.NotEmpty(t, homeset)
	e := &Event{
		Name:          "VEVENT",
		Summary:       "test-event",
		Uid:           modificateUID,
		DateTimeStart: time.Now(),
		DateTimeEnd:   time.Now(),
		Attendees:     []string{testEmail2, "likh.lyudmila1@yandex.ru"},
		Organizer:     testEmail1,
	}
	event, err := GetEvent(e)
	assert.NoError(t, err)

	err = CreateEvent(ctx, client, homeset, modCalendarName, event)
	assert.NoError(t, err)
}

func TestAttend_Reply(t *testing.T) {
	var client *caldav.Client
	var principal, homeset string
	var ctx context.Context
	var err error
	_, client, principal, ctx, err = CreateClient(URL, bytes.NewBufferString("ya\niamyandex\n"))
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err = client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)
	assert.NotEmpty(t, homeset)

	resp, err := GetEvents(ctx, client, homeset, inboxCalendarName)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	var uid string
	var r caldav.CalendarObject
	for _, r = range resp {
		uid, err = r.Data.Props.Text(ical.PropUID)
		assert.NoError(t, err)
		if uid == modificateUID {
			break
		}
	}

	eventFileName := r.Path
	eventFileName = eventFileName[len(homeset+inboxCalendarName+"/") : len(eventFileName)-len(".ics")]

	var mods *Modifications = &Modifications{
		PartStat:     "ACCEPTED",
		LastModified: time.Now(),
		DelegateTo:   "",
		CalendarName: modCalendarName,
	}
	err = ModifyAttendance(ctx, client, homeset, inboxCalendarName, modificateUID, eventFileName, mods)
	assert.NoError(t, err)

}
func TestAttend_CheckStatus(t *testing.T) {
	var client *caldav.Client
	var principal, homeset string
	var ctx context.Context
	var err error
	var resp []caldav.CalendarObject

	_, client, principal, ctx, err = CreateClient(URL, bytes.NewBufferString(testCredentials1))
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err = client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)
	assert.NotEmpty(t, homeset)

	resp, err = GetEvents(ctx, client, homeset, modCalendarName)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	for _, r := range resp {
		uid, err := r.Data.Props.Text(ical.PropUID)
		assert.NoError(t, err)
		if uid == modificateUID {
			stat, err := r.Data.Props.Text(ical.ParamParticipationStatus)
			assert.NoError(t, err)
			assert.Equal(t, "ACCEPTED", stat)
			break
		}
	}
}

func TestAttendClearAll(t *testing.T) {
	var client *caldav.Client
	var principal, homeset string
	var ctx context.Context
	var err error
	var resp []caldav.CalendarObject

	_, client, principal, ctx, err = CreateClient(URL, bytes.NewBufferString(testCredentials1))
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err = client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)
	assert.NotEmpty(t, homeset)
	resp, err = GetEvents(ctx, client, homeset, modCalendarName)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	for _, r := range resp {
		err = Delete(ctx, client, r.Path)
		assert.NoError(t, err)
	}
	_, client, principal, ctx, err = CreateClient(URL, bytes.NewBufferString(testCredentials2))
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err = client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)
	assert.NotEmpty(t, homeset)
	resp, err = GetEvents(ctx, client, homeset, modCalendarName)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	for _, r := range resp {
		err = Delete(ctx, client, r.Path)
		assert.NoError(t, err)
	}

	_, client, principal, ctx, err = CreateClient(URL, bytes.NewBufferString(testCredentials1))
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err = client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)
	assert.NotEmpty(t, homeset)
	resp, err = GetEvents(ctx, client, homeset, inboxCalendarName)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	for _, r := range resp {
		err = Delete(ctx, client, r.Path)
		assert.NoError(t, err)
	}
	_, client, principal, ctx, err = CreateClient(URL, bytes.NewBufferString(testCredentials2))
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err = client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)
	assert.NotEmpty(t, homeset)
	resp, err = GetEvents(ctx, client, homeset, inboxCalendarName)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	for _, r := range resp {
		err = Delete(ctx, client, r.Path)
		assert.NoError(t, err)
	}
}

func TestAddAttendee(t *testing.T) {
	var client *caldav.Client
	var principal, homeset string
	var ctx context.Context
	var err error

	_, client, principal, ctx, err = CreateClient(URL, bytes.NewBufferString(testCredentials1))
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err = client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)
	assert.NotEmpty(t, homeset)

	err = PutAttendee(ctx, client, "likh.lyudmila1@yandex.ru", homeset, modCalendarName, modificateUID, modificateUID)
	assert.NoError(t, err)
}

func TestAttendees(t *testing.T) {
	var client *caldav.Client
	var ctx context.Context
	var err error

	_, client, _, ctx, err = CreateClient(URL, bytes.NewBufferString(testCredentials1))
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, ctx)

	comp, err := FindEvent(ctx, client, "/dav.php/calendars/testuser/mod/27518c07-3f45-11ef-a928-80d21df4779b.ics", "27518c07-3f45-11ef-a928-80d21df4779b")
	if err != nil {
		fmt.Printf("1 "+"%s\n", err)
		return
	}
	// props - map [string][]prop
	// prop - name, params, value
	// params- map [string][]string

	attendeeProp := comp.Props.Get(ical.PropAttendee)
	fmt.Printf("attendee prop name: %s\n", attendeeProp.Name)
	fmt.Printf("value type: %v\n", attendeeProp.ValueType())
	addr, err := comp.Props.URI(ical.PropAttendee)
	if err != nil {
		fmt.Printf("2 %s\n", err)
	}
	fmt.Println(addr)

}