package ui

import (
	"context"
	"fmt"
	"os"

	"github.com/emersion/go-webdav/caldav"
	"golang.org/x/term"
)

func ClearLines(num int) {
	for i := 0; i < num; i++ {
		fmt.Print("\033[A")
		fmt.Print("\033[2K")
	}
}

func ColouredLine(str string) {
	fmt.Printf("\u001b[34m%s\u001b[0m\n", str)
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
	ClearLines(3)
	return username, password
}

func StartMenu() {
	ColouredLine("Main menu:")
	for {
		fmt.Println("1. Log in")
		fmt.Println("0. Exit")
		var answer int
		fmt.Scan(&answer)
		ClearLines(3)
		switch answer {
		case 1:
			CalendarMenu(CreateClient())
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
			ListCalendars(client, ctx)
		case 2:
			EventMenu("FIX", ctx) // FIX!!! what to pass?
		case 3:
			//DeleteCalendar(client)
		case 4:
			//ListEvents(client)
		case 5:
			//CreateEvent(client)
		case 6:
			//DeleteEvent(client)
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
