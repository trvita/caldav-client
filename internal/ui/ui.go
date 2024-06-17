package ui

import (
	"context"
	"fmt"
	"log"

	"github.com/emersion/go-webdav/caldav"
	mycal "github.com/trvita/caldav-client/internal/caldav"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("\u001b[31m%s: %s\u001b[0m\n", msg, err)
	}
}

func ClearLines(num int) {
	for i := 0; i < num; i++ {
		fmt.Print("\033[A")
		fmt.Print("\033[2K")
	}
}

func ColouredLine(str string) {
	fmt.Printf("\u001b[34m%s\u001b[0m\n", str)
}

func StartMenu(url string) {
	ColouredLine("Main menu:")
	for {
		fmt.Println("1. Log in")
		fmt.Println("0. Exit")
		var answer int
		fmt.Scan(&answer)
		ClearLines(3)
		switch answer {
		case 1:
			client, ctx := mycal.CreateClient(url)
			ClearLines(3)
			CalendarMenu(client, ctx)
		case 0:
			ColouredLine("Shutting down...")
			return
		}
	}
}

func CalendarMenu(client *caldav.Client, ctx context.Context) {
	principal, err := client.FindCurrentUserPrincipal(ctx)
	FailOnError(err, "Error finding current user principal")
	ColouredLine("Current user: " + principal[1:len(principal)-1])
	for {
		fmt.Println("1. List calendars")
		fmt.Println("2. Goto calendar")
		fmt.Println("3. Create calendar")
		fmt.Println("4. Delete calendar")
		// fmt.Println("5. Update calendar")
		fmt.Println("0. Log out")
		var answer int
		fmt.Scan(&answer)
		ClearLines(6)
		switch answer {
		case 1:
			mycal.ListCalendars(client, ctx)
		case 2:
			EventMenu("FIX", ctx) // FIX!!! what to pass?
		case 3:
			//mycal.DeleteCalendar(client)
		case 4:
			//mycal.ListEvents(client)
		case 5:
			//mycal.CreateEvent(client)
		case 6:
			//mycal.DeleteEvent(client)
		case 0:
			ColouredLine("Logging out...")
			return
		}
	}
}

func EventMenu(calendar string, ctx context.Context) {
	ColouredLine("Current calendar:" + calendar)
	for {
		fmt.Println("1. List calendars")
		fmt.Println("2. Goto calendar")
		fmt.Println("3. Create calendar")
		fmt.Println("4. Delete calendar")
		// fmt.Println("5. Update calendar")
		fmt.Println("0. Back to calendar menu")
		var answer int
		fmt.Scan(&answer)
		ClearLines(6)
		switch answer {
		case 1:
		case 2:
		case 3:
		case 0:
			ColouredLine("Returning to calendar menu...")
			return
		}
	}
}
