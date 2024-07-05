package caldav

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
	"github.com/teambition/rrule-go"
	"golang.org/x/term"
)

type Event struct {
	Name          string
	Summary       string
	Uid           string
	DateTimeStart time.Time
	DateTimeEnd   time.Time
	Attendees     []string
	Organizer     string
}

type ReccurentEvent struct {
	Name          string
	Summary       string
	Uid           string
	DateTimeStart time.Time
	DateTimeUntil time.Time
	Frequency     int
	Count         int
	Interval      int
	ByDay         []int
	ByMonthDay    []int
	ByYearDay     []int
	ByMonth       []int
	ByWeekNo      []int
	ByHour        []int
	BySetPos      []int
	Attendees     []string
	Organizer     string
}

func ExtractNameFromEmail(email string) string {
	emailParts := strings.Split(email, "@")
	if len(emailParts) != 2 {
		return ""
	}
	return emailParts[0]
}

func GetCredentials(r io.Reader) (string, string, error) {
	reader := bufio.NewReader(r)
	fmt.Print("username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	username = strings.TrimSpace(username)

	fmt.Print("password: ")
	var password string
	if r == os.Stdin {
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", "", err
		}
		password = string(bytePassword)
		fmt.Println()
	} else {
		password, err = reader.ReadString('\n')
		if err != nil {
			return "", "", err
		}
		password = strings.TrimSpace(password)
	}
	return username, password, nil
}

func CreateClient(url string, r io.Reader) (webdav.HTTPClient, *caldav.Client, string, context.Context, error) {
	username, password, err := GetCredentials(r)

	if err != nil {
		return nil, nil, "", nil, err
	}

	httpClient := webdav.HTTPClientWithBasicAuth(&http.Client{}, username, password)
	client, err := caldav.NewClient(httpClient, url)
	if err != nil {
		return nil, nil, "", nil, err
	}

	ctx := context.Background()
	principal, err := client.FindCurrentUserPrincipal(ctx)
	if err != nil {
		return nil, nil, "", nil, err
	}
	return httpClient, client, principal, ctx, nil
}

func ListCalendars(ctx context.Context, client *caldav.Client, homeset string) error {
	calendars, err := client.FindCalendars(ctx, homeset)
	if err != nil {
		return err
	}
	// maybe return calendars and not use print in caldav.go
	for _, calendar := range calendars {
		fmt.Printf("Calendar: %s\n", calendar.Name)
	}
	return nil
}

func CreateCalendar(ctx context.Context, httpClient webdav.HTTPClient, url, homeset, calendarName, description string) error {
	reqBody := fmt.Sprintf(`
	<C:mkcalendar xmlns:D='DAV:' xmlns:C='urn:ietf:params:xml:ns:caldav'>
			<D:set>
				<D:prop>
					<D:displayname>%s</D:displayname>
					<C:calendar-description>%s</C:calendar-description>
				</D:prop>
			</D:set>
		</C:mkcalendar>`, calendarName, description)
	calURL := url + homeset[8:] + calendarName
	req, err := http.NewRequest("MKCALENDAR", calURL, bytes.NewBufferString(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/xml")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return err
	}
	return nil
}

func FindCalendar(ctx context.Context, client *caldav.Client, homeset, calendarName string) error {
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

func ListEvents(ctx context.Context, client *caldav.Client, homeset, calendarName string) error {
	query := &caldav.CalendarQuery{
		CompRequest: caldav.CalendarCompRequest{
			Name:     "VCALENDAR",
			AllProps: true,
			Comps: []caldav.CalendarCompRequest{{
				Name:     "VEVENT",
				AllProps: true,
			}},
		},
		CompFilter: caldav.CompFilter{
			Name: "VCALENDAR",
			Comps: []caldav.CompFilter{{
				Name: "VEVENT",
			}},
		},
	}

	calendarURL := homeset + calendarName
	resp, err := client.QueryCalendar(ctx, calendarURL, query)
	if err != nil {
		return fmt.Errorf("error getting calendar query: %v", err)
	}
	if len(resp) == 0 {
		return fmt.Errorf("no events found")
	}
	// maybe return props and not use print in caldav.go
	for _, calendarObject := range resp {
		for _, event := range calendarObject.Data.Children {
			for _, prop := range event.Props {
				for _, p := range prop {
					fmt.Printf("%s: %s\n", p.Name, p.Value)
				}
			}
			fmt.Println()
		}
	}
	return nil
}

func ListTodos(ctx context.Context, client *caldav.Client, homeset, calendarName string) error {
	query := &caldav.CalendarQuery{
		CompRequest: caldav.CalendarCompRequest{
			Name: "VCALENDAR",
			Comps: []caldav.CalendarCompRequest{{
				Name: "VTODO",
			}},
		},
		CompFilter: caldav.CompFilter{
			Name: "VCALENDAR",
			Comps: []caldav.CompFilter{{
				Name: "VTODO",
			}},
		},
	}

	calendarURL := homeset + calendarName
	resp, err := client.QueryCalendar(ctx, calendarURL, query)
	if err != nil {
		return fmt.Errorf("error getting calendar query: %v", err)
	}
	if len(resp) == 0 {
		return fmt.Errorf("no todos found")
	}
	// maybe return props and not use print in caldav.go
	for _, calendarObject := range resp {
		for _, event := range calendarObject.Data.Children {
			for _, prop := range event.Props {
				for _, p := range prop {
					fmt.Printf("%s: %s\n", p.Name, p.Value)
				}
			}
			fmt.Println()
		}
	}
	return nil
}

func GetEvent(newEvent *Event) *ical.Event {
	event := ical.NewEvent()
	event.Name = newEvent.Name
	event.Props.SetText(ical.PropUID, newEvent.Uid)
	event.Props.SetText(ical.PropSummary, newEvent.Summary)
	event.Props.SetDateTime(ical.PropDateTimeStamp, time.Now().UTC())
	event.Props.SetDateTime(ical.PropDateTimeStart, newEvent.DateTimeStart)
	event.Props.SetDateTime(ical.PropDateTimeEnd, newEvent.DateTimeEnd)
	for _, attendee := range newEvent.Attendees {
		prop := ical.NewProp(ical.PropAttendee)
		prop.Params.Add(ical.ParamParticipationStatus, "NEEDS-ACTION")
		prop.Params.Add(ical.ParamRole, "REQ-PARTICIPANT")
		prop.Value = "mailto:" + attendee
		event.Props.Add(prop)
	}
	if newEvent.Attendees != nil {
		propOrg := ical.NewProp(ical.PropOrganizer)
		propOrg.Value = "mailto:" + newEvent.Organizer
		event.Props.Add(propOrg)
	}
	return event
}

func GetTodo(newEvent *Event) *ical.Event {
	event := ical.NewEvent()
	event.Name = newEvent.Name
	event.Props.SetText(ical.PropUID, newEvent.Uid)
	event.Props.SetText(ical.PropSummary, newEvent.Summary)
	event.Props.SetDateTime(ical.PropDateTimeStamp, time.Now().UTC())
	event.Props.SetDateTime(ical.PropLastModified, time.Now().UTC()) // maybe server puts it
	event.Props.SetDateTime(ical.PropDateTimeStart, newEvent.DateTimeStart)
	event.Props.SetDateTime(ical.PropDue, newEvent.DateTimeEnd)
	event.Props.SetText(ical.PropStatus, "NEEDS-ACTION")
	return event
}

func GetRecurrentEvent(newRecEvent *ReccurentEvent) *ical.Event {
	event := ical.NewEvent()
	event.Name = newRecEvent.Name
	event.Props.SetText(ical.PropUID, newRecEvent.Uid)
	event.Props.SetText(ical.PropSummary, newRecEvent.Summary)
	event.Props.SetDateTime(ical.PropDateTimeStamp, time.Now().UTC())
	event.Props.SetDateTime(ical.PropDateTimeStart, newRecEvent.DateTimeStart)

	event.Props.SetRecurrenceRule(&rrule.ROption{
		Freq: rrule.Frequency(newRecEvent.Frequency),
		//Dtstart:    newRecEvent.DateTimeStart,
		Interval:   newRecEvent.Interval,
		Wkst:       rrule.MO,
		Count:      newRecEvent.Count,
		Until:      newRecEvent.DateTimeUntil,
		Bysetpos:   newRecEvent.BySetPos,
		Bymonth:    newRecEvent.ByMonth,
		Bymonthday: newRecEvent.ByMonthDay,
		Byyearday:  newRecEvent.ByYearDay,
		Byweekno:   newRecEvent.ByWeekNo,
		Byweekday:  []rrule.Weekday{},
		Byhour:     newRecEvent.ByHour,
		Byminute:   []int{},
		Bysecond:   []int{},
		Byeaster:   []int{},
	})

	for _, attendee := range newRecEvent.Attendees {
		prop := ical.NewProp(ical.PropAttendee)
		prop.Params.Add(ical.ParamParticipationStatus, "NEEDS-ACTION")
		prop.Params.Add(ical.ParamRole, "REQ-PARTICIPANT")
		prop.Value = "mailto:" + attendee
		event.Props.Add(prop)
	}
	if newRecEvent.Attendees != nil {
		propOrg := ical.NewProp(ical.PropOrganizer)
		propOrg.Value = "mailto:" + newRecEvent.Organizer
		event.Props.Add(propOrg)
	}
	return event
}

func CreateEvent(ctx context.Context, client *caldav.Client, homeset string, calendarName string, event *ical.Event) error {
	calendar := ical.NewCalendar()
	calendar.Props.SetText(ical.PropVersion, "2.0")
	calendar.Props.SetText(ical.PropProductID, "-//trvita//EN")
	calendar.Props.SetText(ical.PropCalendarScale, "GREGORIAN")

	timezone := ical.NewComponent(ical.CompTimezone)
	timezone.Props.SetText(ical.ParamTimezoneID, "Asia/Krasnoyarsk")

	standard := ical.NewComponent(ical.CompTimezoneStandard)
	standard.Props.SetDateTime(ical.PropDateTimeStart, time.Date(1970, 01, 01, 0, 0, 0, 0, time.Now().Location()))
	standard.Props.SetText(ical.PropTimezoneOffsetFrom, "+0700")
	standard.Props.SetText(ical.PropTimezoneOffsetTo, "+0700")
	timezone.Children = append(timezone.Children, standard)
	calendar.Children = append(calendar.Children, timezone)

	calendar.Children = append(calendar.Children, event.Component)
	eventUID, err := event.Props.Text(ical.PropUID)
	if err != nil {
		return err
	}
	eventURL := homeset + calendarName + "/" + eventUID + ".ics"
	_, err = client.PutCalendarObject(ctx, eventURL, calendar)
	if err != nil {
		return err
	}
	return nil
}

func Delete(ctx context.Context, client *caldav.Client, path string) error {
	err := client.RemoveAll(ctx, path)
	if err != nil {
		return err
	}
	return nil
}
