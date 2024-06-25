package ui

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/emersion/go-webdav/caldav"
	"github.com/google/uuid"
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
	fmt.Printf("\u001b[34m%s\u001b[0m", str)
}

func GetString(message string) string {
	var str string
	fmt.Print(message)
	fmt.Scan(&str)
	return str
}

func GetEvent() (string, string, time.Time, time.Time) {
	var summary, startDate, startTime, endDate, endTime string
	var startDateTime, endDateTime time.Time
	var err error
	summary = GetString("Enter event summary: ")
	uid, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("could not generate UUID: %v", err)
	}
	for {
		startDate = GetString("Enter event start date (YYYY.MM.DD): ")
		startTime = GetString("Enter event start time (HH.MM.SS): ")

		startDateTime, err = time.Parse("2006.01.02 15.04.05", startDate+" "+startTime)
		if err != nil {
			fmt.Println("invalid start date/time format")
			continue
		}
		break
	}
	for {
		endDate = GetString("Enter event end date (YYYY.MM.DD): ")
		endTime = GetString("Enter event end time (HH.MM.SS): ")

		endDateTime, err = time.Parse("2006.01.02 15.04.05", endDate+" "+endTime)
		if err != nil {
			fmt.Println("invalid end date/time format")
			continue
		}
		break
	}
	return summary, uid.String(), startDateTime, endDateTime
}

func StartMenu(url string) {
	ColouredLine("Main menu:\n")
	for {
		fmt.Println("1. Log in")
		fmt.Println("0. Exit")
		var answer int
		fmt.Scan(&answer)
		ClearLines(3)
		switch answer {
		case 1:
			var client *caldav.Client
			var principal string
			var ctx context.Context
			var err error
			for {
				client, principal, ctx, err = mycal.CreateClient(url, os.Stdin)
				ClearLines(2)
				if err == nil {
					break
				}
			}
			CalendarMenu(client, principal, ctx)
		case 0:
			ColouredLine("Shutting down...\n")
			return
		}
	}
}

func CalendarMenu(client *caldav.Client, principal string, ctx context.Context) {
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	FailOnError(err, "Error finding calendar homeset")
	ColouredLine("Current user: " + principal[1:len(principal)-1] + "\n")
	for {
		fmt.Println("1. List calendars")
		fmt.Println("2. Goto calendar")
		fmt.Println("3. Create calendar")
		fmt.Println("0. Log out")
		var answer int
		fmt.Scan(&answer)
		ClearLines(5)
		switch answer {
		case 1:
			mycal.ListCalendars(ctx, client, homeset)
		case 2:
			calendarName := GetString("Enter new calendar name:")
			err := mycal.FindCalendar(ctx, client, homeset, calendarName)
			if err != nil {
				fmt.Printf("%s\n", err)
				break
			}
			EventMenu(ctx, client, homeset, calendarName)
		case 3:
			calendarName := GetString("Enter calendar name:")
			summary, uid, startDateTime, endDateTime := GetEvent()
			mycal.CreateCalendar(ctx, client, homeset, calendarName, summary, uid, startDateTime, endDateTime)
		case 0:
			ColouredLine("Logging out...\n")
			return
		}
	}
}

func EventMenu(ctx context.Context, client *caldav.Client, homeset string, calendar string) {
	ColouredLine("Current calendar:" + calendar + "\n")
	for {
		fmt.Println("1. List events")
		fmt.Println("2. Create event")
		fmt.Println("3. Delete event")
		fmt.Println("0. Back to calendar menu")
		var answer int
		fmt.Scan(&answer)
		ClearLines(5)
		switch answer {
		case 1:
			mycal.ListEvents(ctx, client, homeset, calendar)
		case 2:
			summary, uid, startDateTime, endDateTime := GetEvent()
			event := mycal.GetEvent(summary, uid, startDateTime, endDateTime)
			mycal.CreateEvent(ctx, client, homeset, calendar, event)
		case 3:
			eventUID := GetString("Enter event UID: ")
			mycal.DeleteEvent(ctx, client, homeset, calendar, eventUID)
		case 0:
			ColouredLine("Returning to calendar menu...\n")
			return
		}
	}
}
