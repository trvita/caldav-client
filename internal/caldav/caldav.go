package caldav

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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
	username = username[:len(username)-1]

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
		password = password[:len(password)-1]
	}
	return username, password
}

func CreateClient(url string, r io.Reader) (*caldav.Client, context.Context, error) {
	username, password := GetCredentials(r)
	httpClient := webdav.HTTPClientWithBasicAuth(&http.Client{}, username, password)
	client, err := caldav.NewClient(httpClient, url)
	if err != nil {
		return nil, nil, err
	}
	return client, context.Background(), nil
}

func ListCalendars(client *caldav.Client, ctx context.Context) {
	principal, err := client.FindCurrentUserPrincipal(ctx)
	FailOnError(err, "Error finding current user principal")
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	FailOnError(err, "Error finding calendar home set")
	calendars, err := client.FindCalendars(ctx, homeset)
	FailOnError(err, "Error fetching calendars")
	for _, calendar := range calendars {
		fmt.Printf("Calendar: %s\n", calendar.Name)
	}
}

func CreateCalendar(client *caldav.Client, ctx context.Context, calendarName string) {
	principal, err := client.FindCurrentUserPrincipal(ctx)
	FailOnError(err, "Error finding current user principal")
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	FailOnError(err, "Error finding calendar home set")
	newCalendar := ical.NewCalendar()
	// newCalendar.Props.SetDateTime(ical.PropTimezoneID, time.Now())
	// newCalendar.Props.SetText(ical.PropVersion, "2.0")
	// newCalendar.Props.SetText(ical.PropProductID, "-//trvita//caldav-client//EN")
	// newCalendar.Props.SetText(ical.PropCalendarScale, "GREGORIAN")

	calendarPath := homeset + "/" + calendarName + ".ics"
	_, err = client.PutCalendarObject(ctx, calendarPath, newCalendar)
	FailOnError(err, "Error creating calendar")
}

func ListEvents(client *caldav.Client, ctx context.Context, calendarName string)  {}
func FindEvent(client *caldav.Client, ctx context.Context, calendarName string)   {}
func CreateEvent(client *caldav.Client, ctx context.Context, calendarName string) {}
func DeleteEvent(client *caldav.Client, ctx context.Context, calendarName string) {}
