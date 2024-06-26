package caldav

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
	"golang.org/x/term"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("\u001b[31m%s: %s\u001b[0m\n", msg, err)
	}
}

func GetCredentials(r io.Reader) (string, string) {
	reader := bufio.NewReader(r)
	fmt.Print("username: ")
	username, err := reader.ReadString('\n')
	FailOnError(err, "Error reading username")
	username = strings.TrimSpace(username)

	fmt.Print("password: ")
	var password string
	if r == os.Stdin {
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		FailOnError(err, "Error reading password")
		password = string(bytePassword)
		fmt.Println()
	} else {
		password, err = reader.ReadString('\n')
		FailOnError(err, "Error reading password")
		password = strings.TrimSpace(password)
	}
	return username, password
}

func CreateClient(url string, r io.Reader) (*caldav.Client, string, context.Context, error) {
	username, password := GetCredentials(r)
	httpClient := webdav.HTTPClientWithBasicAuth(&http.Client{}, username, password)
	client, err := caldav.NewClient(httpClient, url)
	if err != nil {
		return nil, "", nil, err
	}

	ctx := context.Background()
	principal, err := client.FindCurrentUserPrincipal(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	return client, principal, ctx, nil
}

func ListCalendars(ctx context.Context, client *caldav.Client, homeset string) {
	calendars, err := client.FindCalendars(ctx, homeset)
	FailOnError(err, "Error fetching calendars")
	for _, calendar := range calendars {
		fmt.Printf("Calendar: %s\n", calendar.Name)
	}
}

func CreateCalendar(ctx context.Context, client *caldav.Client, homeset string, calendarName string, summary string, uid string, startDateTime time.Time, endDateTime time.Time) {
	calendar := ical.NewCalendar()
	calendar.Props.SetText(ical.PropVersion, "2.0")
	calendar.Props.SetText(ical.PropProductID, "-//trvita//EN")
	calendar.Props.SetText(ical.PropCalendarScale, "GREGORIAN")

	event := GetEvent(summary, uid, startDateTime, endDateTime)

	calendar.Children = append(calendar.Children, event.Component)
	calendarURL := homeset + calendarName + "/"
	_, err := client.PutCalendarObject(ctx, calendarURL, calendar)
	FailOnError(err, "Error putting calendar object")
	fmt.Println("Calendar created")
}

func FindCalendar(ctx context.Context, client *caldav.Client, homeset string, calendarName string) error {
	calendars, err := client.FindCalendars(ctx, homeset)
	if err != nil {
		return err
	}

	for _, calendar := range calendars {
		if calendar.Name == calendarName {
			return nil
		}
	}
	return fmt.Errorf("calendar with name %s not found", calendarName)
}

func ListEvents(ctx context.Context, client *caldav.Client, homeset string, calendarName string) {
	fmt.Printf("Events:\n")
	calendar, err := client.GetCalendarObject(ctx, calendarName)
	FailOnError(err, "Error getting calendar object")
	for _, event := range calendar.Data.Events() {
		summary, err := event.Props.Text("SUMMARY")
		FailOnError(err, "Error reading summary")
		uid, err := event.Props.Text("UID")
		FailOnError(err, "Error reading UID")
		fmt.Printf("Summary: %s,\tUID: %s\n", summary, uid)
	}

}

func GetEvent(summary string, uid string, startDateTime time.Time, endDateTime time.Time) *ical.Event {
	event := ical.NewEvent()
	event.Props.SetText(ical.PropUID, uid)
	event.Props.SetText(ical.PropSummary, summary)
	event.Props.SetDateTime(ical.PropDateTimeStamp, time.Now().UTC())
	event.Props.SetDateTime(ical.PropDateTimeStart, startDateTime)
	event.Props.SetDateTime(ical.PropDateTimeEnd, endDateTime)
	fmt.Println("Event created with UID " + uid)
	return event
}

func CreateEvent(ctx context.Context, client *caldav.Client, homeset string, calendarName string, event *ical.Event) {
	calendar, err := client.GetCalendarObject(ctx, calendarName)
	FailOnError(err, "Error getting calendar object")
	calendar.Data.Component.Children = append(calendar.Data.Component.Children, event.Component)
	var buf strings.Builder
	encoder := ical.NewEncoder(&buf)
	err = encoder.Encode(calendar.Data)
	FailOnError(err, "Error encoding calendar")
	_, err = client.PutCalendarObject(ctx, calendarName, calendar.Data)
	FailOnError(err, "Error putting calendar object")
}

func DeleteEvent(ctx context.Context, client *caldav.Client, homeset string, calendarName string, eventUID string) {
	calendar, err := client.GetCalendarObject(ctx, calendarName)
	FailOnError(err, "Error getting calendar object")
	var updatedEvents []*ical.Component
	for _, component := range calendar.Data.Component.Children {
		if component.Name == ical.CompEvent {
			uid, err := component.Props.Text(ical.PropUID)
			FailOnError(err, "Error reading UID")
			if uid == eventUID {
				continue
			}
		}
		updatedEvents = append(updatedEvents, component)
	}
	if len(updatedEvents) == 0 {
		fmt.Println("Cannot delete the event as it would leave the calendar empty.") // add delete calendar call if implemented
		return
	}

	calendar.Data.Component.Children = updatedEvents

	var buf strings.Builder
	encoder := ical.NewEncoder(&buf)
	err = encoder.Encode(calendar.Data)
	FailOnError(err, "Error encoding calendar")
	_, err = client.PutCalendarObject(ctx, calendarName, calendar.Data)
	FailOnError(err, "Error putting calendar object")

	fmt.Println("Event deleted")
}
