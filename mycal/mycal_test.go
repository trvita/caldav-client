package mycal

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	webdav "github.com/trvita/caldav-client-yandex"
	"github.com/trvita/caldav-client-yandex/caldav"
	"github.com/trvita/go-ical"
)

var URL = "http://127.0.0.1:90/dav.php"
var credentials = 0
var email = 1

var users = map[string][]string{
	"user1": {"testuser\ntestpassword\n", "some-mail@mail.com"},
	"user2": {"usertest\npasswordtest\n", "mail-some@mail.com"},
	"user3": {"tuserest\ntasswordpest\n", "maso-meil@mail.com"},
}

var calendars = []string{
	"cal-new", "default", "cal-empty", "wrong", "inbox",
}

var calendarsAttend = []string{
	"some-cal",
}

var uids = []string{
	"valid", "invalid", "modificate", "modificate-not-share", "shared",
}

func setupClient(t *testing.T, user string) (webdav.HTTPClient, *caldav.Client, string, context.Context) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(users[user][credentials]))
	assert.NoError(t, err)
	assert.NotNil(t, httpClient)
	assert.NotNil(t, client)
	assert.NotEmpty(t, principal)
	assert.NotNil(t, ctx)
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	assert.NoError(t, err)
	assert.NotNil(t, homeset)
	return httpClient, client, homeset, ctx
}

func TestGetCredentials(t *testing.T) {
	input := bytes.NewBufferString(users["user1"][credentials])
	username, password, err := GetCredentials(input)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", username)
	assert.Equal(t, "testpassword", password)
}

func TestCreateClient(t *testing.T) {
	httpClient, client, principal, ctx, err := CreateClient(URL, bytes.NewBufferString(users["user1"][credentials]))
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
	_, client, homeset, ctx := setupClient(t, "user1")

	err := ListCalendars(ctx, client, homeset)
	assert.NoError(t, err)
}

func TestCreateCalendar(t *testing.T) {
	httpClient, client, homeset, ctx := setupClient(t, "user1")

	err := CreateCalendar(ctx, httpClient, URL, homeset, calendars[0], "0")
	assert.NoError(t, err)

	err = Delete(ctx, client, homeset+calendars[0])
	assert.NoError(t, err)
}

func TestFindCalendarCorrect(t *testing.T) {
	_, client, homeset, ctx := setupClient(t, "user1")

	err := FindCalendar(ctx, client, homeset, calendars[1])
	assert.NoError(t, err)
}

func TestFindCalendarFail(t *testing.T) {
	_, client, homeset, ctx := setupClient(t, "user1")

	err := FindCalendar(ctx, client, homeset, calendars[3])
	assert.Error(t, err)
}

func TestGetEvents(t *testing.T) {
	_, client, homeset, ctx := setupClient(t, "user1")

	resp, err := GetEvents(ctx, client, homeset, calendars[1])
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}

func TestGetEventsFail(t *testing.T) {
	_, client, homeset, ctx := setupClient(t, "user1")

	resp, err := GetEvents(ctx, client, homeset, calendars[2])
	assert.Error(t, err)
	assert.Empty(t, resp)
}

func TestGetTodos(t *testing.T) {
	_, client, homeset, ctx := setupClient(t, "user1")

	resp, err := GetEvents(ctx, client, homeset, calendars[1])
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}

func TestGetTodosFail(t *testing.T) {
	_, client, homeset, ctx := setupClient(t, "user1")

	resp, err := GetEvents(ctx, client, homeset, calendars[2])
	assert.Error(t, err)
	assert.Empty(t, resp)
}

func TestCreateEvent(t *testing.T) {
	_, client, homeset, ctx := setupClient(t, "user1")

	e := &Event{
		Name:          "VEVENT",
		Summary:       "event",
		Uid:           uids[0],
		DateTimeStart: time.Now(),
		DateTimeEnd:   time.Now(),
		Attendees:     nil,
		Organizer:     "",
	}
	event, err := GetEvent(e)
	assert.NoError(t, err)
	err = CreateEvent(ctx, client, homeset, calendars[1], event)
	assert.NoError(t, err)
}

func TestCreateEventFail(t *testing.T) {
	_, client, homeset, ctx := setupClient(t, "user1")

	e := &Event{}
	event, err := GetEvent(e)
	assert.NoError(t, err)
	err = CreateEvent(ctx, client, homeset, calendars[1], event)
	assert.Error(t, err)
}
func TestDeleteEventFail(t *testing.T) {
	_, client, homeset, ctx := setupClient(t, "user1")

	err := Delete(ctx, client, homeset+calendars[1]+"/"+uids[1]+".ics")
	assert.Error(t, err)
}
func TestDeleteCalendarFail(t *testing.T) {
	_, client, homeset, ctx := setupClient(t, "user1")

	err := Delete(ctx, client, homeset+calendars[3])
	assert.Error(t, err)
}

func TestAttend_Setup(t *testing.T) {
	for user := range users {
		httpClient, _, homeset, ctx := setupClient(t, user)
		for i := range calendarsAttend {
			err := CreateCalendar(ctx, httpClient, URL, homeset, calendarsAttend[i], "")
			assert.NoError(t, err)
		}
	}
}

func TestAttend_Clear(t *testing.T) {
	for user := range users {
		_, client, homeset, ctx := setupClient(t, user)
		for i := range calendarsAttend {
			Delete(ctx, client, homeset+calendarsAttend[i])
		}
		// also clear inbox
		resp, err := GetEvents(ctx, client, homeset, calendars[4])
		if err == nil {
			for _, r := range resp {
				err = Delete(ctx, client, r.Path)
				assert.NoError(t, err)
			}
		}
	}
}

func TestAttend_Invite(t *testing.T) {
	currentUser := "user1"
	_, client, homeset, ctx := setupClient(t, currentUser)

	assert.NotEmpty(t, homeset)
	e := &Event{
		Name:          "VEVENT",
		Summary:       "shared-event",
		Uid:           uids[4],
		DateTimeStart: time.Now(),
		DateTimeEnd:   time.Now(),
		Attendees:     []string{users["user3"][email], users["user2"][email]},
		Organizer:     users[currentUser][email],
	}
	event, err := GetEvent(e)
	assert.NoError(t, err)

	err = CreateEvent(ctx, client, homeset, calendarsAttend[0], event)
	assert.NoError(t, err)
}

func TestAttend_Reply(t *testing.T) {
	currentUser := "user3"
	currentUID := uids[4]
	_, client, homeset, ctx := setupClient(t, currentUser)

	resp, err := GetEvents(ctx, client, homeset, calendars[4])
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	var uid string
	var r caldav.CalendarObject
	for _, r = range resp {
		uid, err = r.Data.Props.Text(ical.PropUID)
		assert.NoError(t, err)
		if uid == currentUID {
			break
		}
	}
	assert.NotEmpty(t, r)
	eventFileName := r.Path
	eventFileName = eventFileName[len(homeset+calendars[4]+"/") : len(eventFileName)-len(".ics")]

	var mods *Modifications = &Modifications{
		PartStat:     "ACCEPTED",
		LastModified: time.Now(),
		DelegateTo:   "",
		CalendarName: calendars[1],
		Email:        users[currentUser][email],
	}
	err = ModifyAttendance(ctx, client, homeset, calendars[4], currentUID, eventFileName, mods)
	assert.NoError(t, err)
}
func TestAttend_Check(t *testing.T) {
	currentUser := "user1"
	currentUID := uids[4]
	currentCalendar := calendars[4]
	_, client, homeset, ctx := setupClient(t, currentUser)

	resp, err := GetByUid(ctx, client, homeset, currentCalendar, currentUID)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	status := resp[0].Data.Children[0].Props.Get(ical.PropAttendee).Params.Get(ical.ParamParticipationStatus)
	assert.Equal(t, "ACCEPTED", status)
	// // at the moment there should be only one event, so there is no reason to go in loop
	// for _, r := range resp {
	// 	for _, event := range r.Data.Children {
	// 		att := event.Props.Get(ical.PropAttendee)
	// 		status := att.Params.Get(ical.ParamParticipationStatus)
	// 		fmt.Println(status)
	// 	}
	// }
}

func TestAttend_notShared(t *testing.T) {
	TestAttend_Clear(t)
	TestAttend_Setup(t)
	TestAttend_Invite(t)
	TestAttend_Reply(t)
	TestAttend_Check(t)
	TestAttend_Clear(t)
}

func TestAddAttendee(t *testing.T) {
	currentUser := "user1"
	_, client, homeset, ctx := setupClient(t, currentUser)

	err := PutAttendee(ctx, client, users["user3"][email], homeset, calendars[5], uids[2], uids[2])
	assert.NoError(t, err)
}

func TestPutEvent(t *testing.T) {
	currentUser := "user1"
	_, client, homeset, ctx := setupClient(t, currentUser)

	resp, err := GetEvents(ctx, client, homeset, calendars[4])
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	var uid string
	var r caldav.CalendarObject
	for _, r = range resp {
		uid, err = r.Data.Props.Text(ical.PropUID)
		assert.NoError(t, err)
		if uid == uids[2] {
			break
		}
	}

	eventFileName := r.Path
	eventFileName = eventFileName[len(homeset+calendars[4]+"/") : len(eventFileName)-len(".ics")]

	var mods *Modifications = &Modifications{
		PartStat:     "",
		LastModified: time.Now(),
		DelegateTo:   "",
		CalendarName: calendars[5],
		Email:        "",
	}
	err = ModifyAttendance(ctx, client, homeset, calendars[4], uids[2], eventFileName, mods)
	assert.NoError(t, err)

}
