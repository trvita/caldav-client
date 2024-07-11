package menu

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	webdav "github.com/trvita/caldav-client-yandex"
	"github.com/trvita/caldav-client-yandex/caldav"
	
	"github.com/trvita/caldav-client/mycal"
	"github.com/trvita/caldav-client/input"
	"github.com/trvita/go-ical"
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

func GetUsernameBaikal(homeset string) string {
	if homeset == "" {
		return ""
	}
	startMarker := URLstart
	startPos := strings.Index(homeset, startMarker)
	if startPos == -1 {
		startPos = 0
	} else {
		startPos += len(startMarker)
	}
	username := homeset[startPos:(len(homeset) - 1)]

	return username
}

func GetModifications(r io.Reader) (*mycal.Modifications, error) {
	var partstat, delegateto, calendarName, email string
	var err error
	var answer byte

	email, err = input.InputString(r, "Enter your email: ")
	if err != nil {
		return nil, err
	}
	fmt.Println("Accept, decline, delegate event? [y, n, d]")
	fmt.Scan(&answer)
	switch answer {
	case 'y':
		partstat = "ACCEPTED"
		calendarName, err = input.InputString(r, "Enter which calendar event goes to: ")
		if err != nil {
			return nil, err
		}
	case 'n':
		partstat = "DECLINED"
	case 'd':
		partstat = string(ical.ParamDelegatedTo)
		delegateto, err = input.InputString(r, "Enter who to delegate: ")
		if err != nil {
			return nil, err
		}
	}
	return &mycal.Modifications{
		Email:        email,
		PartStat:     partstat,
		DelegateTo:   "mailto:" + delegateto,
		CalendarName: calendarName,
		LastModified: time.Now(),
	}, nil
}

func PrintEvents(resp []caldav.CalendarObject) {
	for _, calendarObject := range resp {
		fmt.Printf("path: %s\n", calendarObject.Path)
		for _, event := range calendarObject.Data.Children {
			for _, prop := range event.Props {
				for _, p := range prop {
					fmt.Printf("%s: %s\n", p.Name, p.Value)
				}
			}
			fmt.Println()
		}
	}
}
func StartMenu(url string, r io.Reader) {
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
				httpClient, client, principal, ctx, err = mycal.CreateClient(url, r)
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

			err = CalendarMenu(httpClient, client, principal, ctx, r)
			if err != nil {
				RedLine(err)
			}
		case 0:
			BlueLine("Shutting down...\n")
			return
		}
	}
}

func CalendarMenu(httpClient webdav.HTTPClient, client *caldav.Client, principal string, ctx context.Context, r io.Reader) error {
	homeset, err := client.FindCalendarHomeSet(ctx, principal)
	if err != nil {
		RedLine(err)
		return err
	}
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
			calendarName, err := input.InputString(r, "Enter calendar name to go to:")
			if err != nil {
				return err
			}
			err = mycal.FindCalendar(ctx, client, homeset, calendarName)
			if err != nil {
				RedLine(err)
				break
			}
			EventMenu(ctx, client, homeset, calendarName, r)
		case 3:
			calendarName, err := input.InputString(r, "Enter new calendar name: ")
			if err != nil {
				return err
			}
			description, err := input.InputString(r, "Enter new calendar description: ")
			if err != nil {
				return err
			}
			err = mycal.CreateCalendar(ctx, httpClient, URL, homeset, calendarName, description)
			if err != nil {
				return err
			}
			BlueLine("Calendar " + calendarName + " created\n")
		case 4:
			err := InboxMenu(ctx, client, homeset, "inbox", r)
			if err != nil {
				return err
			}
		case 5:
			calendarName, err := input.InputString(r, "Enter calendar name to delete: ")
			if err != nil {
				return err
			}
			err = mycal.Delete(ctx, client, homeset+calendarName)
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

func EventMenu(ctx context.Context, client *caldav.Client, homeset string, calendarName string, r io.Reader) {
	BlueLine("Current calendar: " + calendarName + "\n")
	for {
		fmt.Println("1. List events")
		fmt.Println("2. Create event")
		fmt.Println("3. Create recurrent event")
		fmt.Println("4. Find events by time range")
		fmt.Println("5. Delete event")
		fmt.Println("0. Back to calendar menu")
		var answer int
		fmt.Scan(&answer)
		switch answer {
		// list events
		case 1:
			BlueLine(calendarName + " EVENTS:\n")
			resp, err := mycal.GetEvents(ctx, client, homeset, calendarName)
			if err != nil {
				RedLine(err)
			}
			PrintEvents(resp)
			BlueLine(calendarName + " TODOS:\n")
			resp, err = mycal.ListTodos(ctx, client, homeset, calendarName)
			if err != nil {
				RedLine(err)
			}
			PrintEvents(resp)
			// create event or todo
		case 2:
			newEvent, err := input.InputEvent(r)
			if err != nil {
				RedLine(err)
				break
			}
			switch newEvent.Name {
			case "VTODO":
				todo, err := mycal.GetTodo(newEvent)
				if err != nil {
					RedLine(err)
					break
				}
				err = mycal.CreateEvent(ctx, client, homeset, calendarName, todo)
				if err != nil {
					RedLine(err)
					break
				}
				BlueLine("Todo created\n")
			case "VEVENT":
				event, err := mycal.GetEvent(newEvent)
				if err != nil {
					RedLine(err)
					break
				}
				err = mycal.CreateEvent(ctx, client, homeset, calendarName, event)
				if err != nil {
					RedLine(err)
					break
				}
				BlueLine("Event created\n")
			}
		case 3:
			newRecEvent, err := input.InputRecurrentEvent(r)
			if err != nil {
				RedLine(err)
				break
			}
			recEvent := mycal.GetRecurrentEvent(newRecEvent)
			err = mycal.CreateEvent(ctx, client, homeset, calendarName, recEvent)
			if err != nil {
				RedLine(err)
				break
			}
			BlueLine("Recurrent event created\n")
		case 4:
			// var startDateTime, endDateTime time.Time
			// var err error
			// for {
			// 	startDate := input.InputString("Enter date to find from (YYYY.MM.DD): ")
			//if err != nil {
			// 	return err
			// }

			// 	startTime := input.InputString("Enter time to find from (HH.MM.SS): ")
			// if err != nil {
			// 	return err
			// }

			// 	startDateTime, err = time.Parse("2006.01.02 15.04.05", startDate+" "+startTime)
			// 	if err != nil {
			// 		fmt.Println("invalid start date/time format")
			// 		continue
			// 	}
			// 	break
			// }
			// for {
			// 	endDate := input.InputString("Enter event end date (YYYY.MM.DD): ")
			// if err != nil {
			// 	return err
			// }

			// 	endTime := input.InputString("Enter event end time (HH.MM.SS): ")
			// if err != nil {
			// 	return err
			// }

			// 	endDateTime, err = time.Parse("2006.01.02 15.04.05", endDate+" "+endTime)
			// 	if err != nil {
			// 		fmt.Println("invalid end date/time format")
			// 		continue
			// 	}
			// 	break
			// }
			// err = mycal.FindEvents(ctx, client, homeset, calendarName, startDateTime, endDateTime)
			// if err != nil {
			// 	RedLine(err)
			// }
		case 5:
			eventUID, err := input.InputString(r, "Enter event path to delete (without .ics): ")
			if err != nil {
				RedLine(err)
				break
			}
			err = mycal.Delete(ctx, client, homeset+calendarName+"/"+eventUID+".ics")
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

func InboxMenu(ctx context.Context, client *caldav.Client, homeset string, calendarName string, r io.Reader) error {
	fmt.Println("Current calendar: ", calendarName)
	for {
		fmt.Println("1. List events")
		fmt.Println("2. Modify event")
		fmt.Println("3. Accept or decline event")
		fmt.Println("0. Return to calendar menu")

		var answer int
		fmt.Scan(&answer)
		switch answer {
		case 1:
			resp, err := mycal.GetEvents(ctx, client, homeset, "inbox")
			if err != nil {
				return err
			}
			PrintEvents(resp)
		case 2:
			eventUID, err := input.InputString(r, "Enter event UID:  ")
			if err != nil {
				return err
			}
			eventPath, err := input.InputString(r, "Enter path to event: ")
			if err != nil {
				return err
			}
			mods, err := GetModifications(r)
			if err != nil {
				return err
			}
			err = mycal.ModifyAttendance(ctx, client, homeset, calendarName, eventUID, eventPath, mods)
			if err != nil {
				return err
			}
		case 0:
			return nil
		}

	}
}
