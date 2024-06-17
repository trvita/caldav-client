package caldav

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/term"

	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("\u001b[31m%s: %s\u001b[0m\n", msg, err)
	}
}

func GetCredentials() (string, string) {
	var username, password string
	fmt.Print("username: ")
	fmt.Scan(&username)
	fmt.Print("password: ")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	FailOnError(err, "Error reading password")
	password = string(bytePassword)
	fmt.Println()
	return username, password
}

func CreateClient(url string) (*caldav.Client, context.Context) {
	username, password := GetCredentials()
	httpClient := webdav.HTTPClientWithBasicAuth(&http.Client{}, username, password)
	client, err := caldav.NewClient(httpClient, url)
	FailOnError(err, "Error creating client")
	return client, context.Background()
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
func CreateCalendar(client *caldav.Client, username string, ctx context.Context) {
	homeset, err := client.FindCalendarHomeSet(ctx, username)
	FailOnError(err, "Error finding calendar home set")
	calendars, err := client.FindCalendars(ctx, homeset)
	FailOnError(err, "Error fetching calendars")
	for _, calendar := range calendars {
		fmt.Printf("Calendar: %s\n", calendar.Name)
	}
}

// func GetCalendarDetails() (string, string) {

// }
