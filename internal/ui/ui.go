package ui

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/emersion/go-ical"
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
	fmt.Printf("\u001b[34m%s\u001b[0m\n", str)
}

func GetCalendarName() string {
	var calendarName string
	fmt.Print("New calendar name: ")
	fmt.Scan(&calendarName)
	return calendarName
}

func GetEvent() *ical.Event {
	var summary, startDate, startTime, endDate, endTime string
	var startDateTime, endDateTime time.Time
	var err error

	fmt.Print("Enter event summary: ")
	fmt.Scan(&summary)
	for {
		fmt.Print("Enter event start date (YYYY.MM.DD): ")
		fmt.Scan(&startDate)

		fmt.Print("Enter event start time (HH.MM.SS): ")
		fmt.Scan(&startTime)

		fmt.Print("Enter event end date (YYYY.MM.DD): ")
		fmt.Scan(&endDate)

		fmt.Print("Enter event end time (HH.MM.SS): ")
		fmt.Scan(&endTime)

		startDateTime, err = time.Parse("2006.01.02 15.04.05", startDate+" "+startTime)
		if err != nil {
			fmt.Println("invalid start date/time format")
			continue
		}

		endDateTime, err = time.Parse("2006.01.02 15.04.05", endDate+" "+endTime)
		if err != nil {
			fmt.Println("invalid end date/time format")
			continue
		}
		break
	}
	event := ical.NewEvent()

	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("could not generate UUID: %v", err)
	}
	event.Props.SetText(ical.PropUID, id.String())
	event.Props.SetText(ical.PropSummary, summary)
	event.Props.SetDateTime(ical.PropDateTimeStamp, time.Now().UTC())
	event.Props.SetDateTime(ical.PropDateTimeStart, startDateTime)
	event.Props.SetDateTime(ical.PropDateTimeEnd, endDateTime)

	return event
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
			ColouredLine("Shutting down...")
			return
		}
	}
}

func CalendarMenu(client *caldav.Client, principal string, ctx context.Context) {
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	FailOnError(err, "Error finding calendar homeset")
	ColouredLine("Current user: " + principal[1:len(principal)-1])
	for {
		fmt.Println("1. List calendars")
		fmt.Println("2. Goto calendar")
		fmt.Println("3. Create calendar")
		fmt.Println("0. Log out")
		var answer int
		fmt.Scan(&answer)
		ClearLines(6)
		switch answer {
		case 1:
			mycal.ListCalendars(ctx, client, homeset)
		case 2:
			calendarName := GetCalendarName()
			EventMenu(ctx, client, homeset, calendarName)
		case 3:
			calendarName := GetCalendarName()
			initialEvent := GetEvent()
			mycal.CreateCalendar(ctx, client, homeset, calendarName, initialEvent)
		case 0:
			ColouredLine("Logging out...")
			return
		}
	}
}

func EventMenu(ctx context.Context, client *caldav.Client, homeset string, calendar string) {
	ColouredLine("Current calendar:" + calendar)
	for {
		fmt.Println("1. List events")
		fmt.Println("3. Create event")
		fmt.Println("4. Delete event")
		fmt.Println("0. Back to calendar menu")
		var answer int
		fmt.Scan(&answer)
		ClearLines(6)
		switch answer {
		case 1:
			mycal.ListEvents(ctx, client, homeset, calendar)
		case 3:
			event := GetEvent()
			mycal.CreateEvent(ctx, client, homeset, calendar, event)
		case 4:
			mycal.DeleteEvent(ctx, client, homeset, calendar)
		case 0:
			ColouredLine("Returning to calendar menu...")
			return
		}
	}
}
