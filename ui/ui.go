package ui

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
	"github.com/google/uuid"
	mycal "github.com/trvita/caldav-client/caldav"
)

var URL = "http://127.0.0.1:90/dav.php"
var URLstart = "/dav.php/calendars/"

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("\u001b[31m%s: %s\u001b[0m\n", msg, err)
	}
}

func BlueLine(str string) {
	fmt.Printf("\u001b[34m%s\u001b[0m", str)
}

func RedLine(err error) {
	fmt.Printf("\u001b[31m%s\u001b[0m\n", err)
}

func GetString(message string) string {
	var str string
	fmt.Print(message)
	fmt.Scan(&str)
	return str
}

func GetStrings(message string) string {
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	str, _ := reader.ReadString('\n')
	return strings.TrimSpace(str)
}

func GetInt(message string) int {
	var num int
	fmt.Print(message)
	fmt.Scan(&num)
	return num
}

func GetUsernameBaikal(homeset string) string {
	startMarker := URLstart
	startPos := strings.Index(homeset, startMarker)
	startPos += len(startMarker)
	username := homeset[startPos:(len(homeset) - 1)]

	return username
}

func GetEvent() (*mycal.Event, error) {
	var attendees []string
	var organizer, name, startDate, startTime, endDate, endTime string
	var startDateTime, endDateTime time.Time
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	for {
		name = GetString("Enter event type [event, todo]: ")
		if strings.ToUpper(name) == "EVENT" {
			name = "VEVENT"
			break
		}
		if strings.ToUpper(name) == "TODO" {
			name = "VTODO"
			break
		}
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
	for {
		attendee := GetString("Enter attendee email (or 0 to finish): ")
		if attendee == "0" {
			break
		}
		attendees = append(attendees, attendee)
	}
	if attendees != nil {
		organizer = GetString("Enter organizer email: ")
	}
	return &mycal.Event{
		Name:          name,
		Uid:           uid.String(),
		Summary:       GetString("Enter event summary: "),
		DateTimeStart: startDateTime,
		DateTimeEnd:   endDateTime,
		Attendees:     attendees,
		Organizer:     organizer,
	}, nil
}

func GetRecurrentEvent() (string, string, time.Time, string, int, int, int, error) {
	var startDate, startTime, name, uid, freq string
	var startDateTime time.Time
	var num, interval, count, until, ans int

	name = "VEVENT"
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", "", time.Time{}, "", 0, 0, 0, err
	}
	uid = uuid.String()
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
	fmt.Print("Enter number of recurrency rules: ")
	fmt.Scan(&num)
	for num > 0 {
		cont := true
		for cont {
			freq = GetString("Enter frequency [Y, MO, D, H, MI]: ")
			switch strings.ToUpper(freq) {
			case "Y":
				freq = "YEARLY"
				cont = false
			case "MO":
				freq = "MONTHLY"
				cont = false
			case "D":
				freq = "DAILY"
				cont = false
			case "H":
				freq = "HOURLY"
				cont = false
			case "MI":
				freq = "MINUTELY"
				cont = false
			}
		}
		interval = GetInt("Enter interval or 0 to skip: ")
		ans = GetInt("Count, until or skip? [1/2/0]: ")
		switch ans {
		case 1:
			count = GetInt("Enter count: ")

		case 2:
			until = GetInt("Enter until: ")
		}
		num--
	}
	return name, uid, startDateTime, freq, interval, count, until, nil
}

func StartMenu(url string) {
	BlueLine("Main menu:\n")
	for {
		fmt.Println("1. Log in")
		fmt.Println("0. Exit")
		var answer int
		fmt.Scan(&answer)
		switch answer {
		case 1:
			var httpClient webdav.HTTPClient
			var client *caldav.Client
			var principal string
			var ctx context.Context
			var err error
			for {
				httpClient, client, principal, ctx, err = mycal.CreateClient(url, os.Stdin)
				if err == nil {
					break
				}
				BlueLine("Wrong username or password, try again? ([y/n])")
				var ans string
				fmt.Scan(&ans)
				ans = strings.ToLower(ans)
				if ans == "y" {
					continue
				} else if ans == "n" {
					BlueLine("Shutting down...\n")
					return
				}
			}
			err = CalendarMenu(httpClient, client, principal, ctx)
			if err != nil {
				RedLine(err)
			}
		case 0:
			BlueLine("Shutting down...\n")
			return
		}
	}
}

func CalendarMenu(httpClient webdav.HTTPClient, client *caldav.Client, principal string, ctx context.Context) error {
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	FailOnError(err, "Error finding calendar homeset")
	for {
		fmt.Println("1. List calendars")
		fmt.Println("2. Goto calendar")
		fmt.Println("3. Create calendar")
		fmt.Println("4. Check inbox")
		fmt.Println("5. Delete calendar")
		fmt.Println("0. Log out")
		var answer int
		fmt.Scan(&answer)
		switch answer {
		case 1:
			err := mycal.ListCalendars(ctx, client, homeset)
			if err != nil {
				RedLine(err)
			}
		case 2:
			calendarName := GetString("Enter calendar name to go to:")
			err := mycal.FindCalendar(ctx, client, homeset, calendarName)
			if err != nil {
				RedLine(err)
				break
			}
			EventMenu(ctx, client, homeset, calendarName)
		case 3:
			calendarName := GetString("Enter new calendar name: ")
			description := GetStrings("Enter new calendar description: ")
			err := mycal.CreateCalendar(ctx, httpClient, URL, homeset, calendarName, description)
			if err != nil {
				RedLine(err)
				break
			}
			BlueLine("Calendar " + calendarName + " created\n")
		case 4:
			err := mycal.ListEvents(ctx, client, homeset, "inbox")
			if err != nil {
				RedLine(err)
			}
		case 5:
			calendarName := GetString("Enter calendar name to delete: ")
			err := mycal.Delete(ctx, client, homeset+calendarName)
			if err != nil {
				RedLine(err)
				break
			}
			BlueLine("Calendar " + calendarName + " deleted\n")
		case 0:
			BlueLine("Logging out...\n")
			return nil
		}
	}
}

func EventMenu(ctx context.Context, client *caldav.Client, homeset string, calendarName string) {
	BlueLine("Current calendar: " + calendarName + "\n")
	for {
		fmt.Println("1. List events")
		fmt.Println("2. Create event")
		fmt.Println("3. Delete event")
		fmt.Println("0. Back to calendar menu")
		var answer int
		fmt.Scan(&answer)
		switch answer {
		// list events
		case 1:
			BlueLine(calendarName + " EVENTS:\n")
			err := mycal.ListEvents(ctx, client, homeset, calendarName)
			if err != nil {
				RedLine(err)
			}
			BlueLine(calendarName + " TODOS:\n")
			err = mycal.ListTodos(ctx, client, homeset, calendarName)
			if err != nil {
				RedLine(err)
			}
			// create event or todo
		case 2:
			newEvent, err := GetEvent()
			if err != nil {
				RedLine(err)
				break
			}
			switch newEvent.Name {
			case "VTODO":
				todo := mycal.GetTodo(newEvent)
				err = mycal.CreateEvent(ctx, client, homeset, calendarName, todo)
				if err != nil {
					RedLine(err)
					break
				}
				BlueLine("Todo created\n")
			case "VEVENT":
				event := mycal.GetEvent(newEvent)
				err = mycal.CreateEvent(ctx, client, homeset, calendarName, event)
				if err != nil {
					RedLine(err)
					break
				}
				BlueLine("Event created\n")
			}
		case 3:
			eventUID := GetString("Enter event UID: ")
			err := mycal.Delete(ctx, client, homeset+calendarName+"/"+eventUID+".ics")
			if err != nil {
				RedLine(err)
				break
			}
			BlueLine("Event deleted\n")
		// go back
		case 0:
			BlueLine("Returning to calendar menu...\n")
			return
		}
	}
}
