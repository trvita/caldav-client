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

func ListCalendars(ctx context.Context, client *caldav.Client, principal string) {
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	FailOnError(err, "Error finding calendar home set")
	calendars, err := client.FindCalendars(ctx, homeset)
	FailOnError(err, "Error fetching calendars")
	for _, calendar := range calendars {
		fmt.Printf("Calendar: %s\n", calendar.Name)
	}
}

func CreateCalendar(ctx context.Context, client *caldav.Client, principal string, calendarName string, event *ical.Event) {
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	FailOnError(err, "Error finding calendar home set")
	calendar := ical.NewCalendar()
	calendar.Props.SetText(ical.PropVersion, "2.0")
	calendar.Props.SetText(ical.PropProductID, "-//trvita//EN")
	calendar.Props.SetText(ical.PropCalendarScale, "GREGORIAN")

	calendar.Children = append(calendar.Children, event.Component)

	var buf strings.Builder
	encoder := ical.NewEncoder(&buf)
	err = encoder.Encode(calendar)
	FailOnError(err, "error encoding calendar")
	calendarURL := homeset + calendarName + "/"
	_, err = client.PutCalendarObject(ctx, calendarURL, calendar)
	FailOnError(err, "Error putting calendar object")
}

func ListEvents(ctx context.Context, client *caldav.Client, calendarName string)  {}
func FindEvent(ctx context.Context, client *caldav.Client, calendarName string)   {}
func CreateEvent(ctx context.Context, client *caldav.Client, calendarName string) {}
func DeleteEvent(ctx context.Context, client *caldav.Client, calendarName string) {}
