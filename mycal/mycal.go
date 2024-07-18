package mycal

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/teambition/rrule-go"
	webdav "github.com/trvita/caldav-client-yandex"
	"github.com/trvita/caldav-client-yandex/caldav"
	"github.com/trvita/go-ical"
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
	Alarm         *Alarm
}

type Alarm struct {
	Action  string
	Trigger string
	// Description string
	// Duration    time.Time
	// Repeat      int
	// Attendee []string
}

type ReccurentEvent struct {
	Event      *Event
	Frequency  int
	Count      int
	Interval   int
	ByDay      []int
	ByMonthDay []int
	ByYearDay  []int
	ByMonth    []int
	ByWeekNo   []int
	ByHour     []int
	BySetPos   []int
}

type Modifications struct {
	Email        string
	PartStat     string
	LastModified time.Time
	DelegateTo   string
	CalendarName string
}

// tested
func ExtractNameFromEmail(email string) string {
	emailParts := strings.Split(email, "@")
	if len(emailParts) != 2 {
		return ""
	}
	return emailParts[0]
}

// tested
func GetCredentials(r io.Reader) (string, string, error) {
	reader := bufio.NewReader(r)
	if r == os.Stdin {
		fmt.Print("username: ")
	}
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	username = strings.TrimSpace(username)

	var password string
	if r == os.Stdin {
		fmt.Print("password: ")
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

// tested
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

// tested
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

// tested
func CreateCalendar(ctx context.Context, httpClient webdav.HTTPClient, url, homeset, calendarName, description string) error {
	timezone := `
	BEGIN:VTIMEZONE
	TZID:Asia/Krasnoyarsk
	TZURL:https://www.tzurl.org/zoneinfo/Asia/Krasnoyarsk
	BEGIN:STANDARD
	TZNAME:+06
	TZOFFSETFROM:+061126
	TZOFFSETTO:+0600
	DTSTART:19200106T000000
	END:STANDARD
	BEGIN:DAYLIGHT
	TZNAME:+07
	TZOFFSETFROM:+0700
	TZOFFSETTO:+0700
	DTSTART:19910331T020000
	END:DAYLIGHT
	END:VTIMEZONE`
	reqBody := fmt.Sprintf(`
	<C:mkcalendar xmlns:D='DAV:' xmlns:C='urn:ietf:params:xml:ns:caldav'>
			<D:set>
				<D:prop>
					<D:displayname>%s</D:displayname>
					<C:calendar-description>%s</C:calendar-description>
					<C:calendar-timezone>%s</C:calendar-timezone>
				</D:prop>
			</D:set>
		</C:mkcalendar>`, calendarName, description, timezone)
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

// tested
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

// tested
func GetEvents(ctx context.Context, client *caldav.Client, homeset, calendarName string) ([]caldav.CalendarObject, error) {
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
		return nil, fmt.Errorf("error getting calendar query: %v", err)
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("no events found")
	}
	return resp, nil
}

func GetByUid(ctx context.Context, client *caldav.Client, homeset, calendarName, uid string) ([]caldav.CalendarObject, error) {
	query := &caldav.CalendarQuery{
		CompRequest: caldav.CalendarCompRequest{
			Name: "VCALENDAR",
			Comps: []caldav.CalendarCompRequest{{
				Name:     "VEVENT",
				AllProps: true,
			}},
		},
		CompFilter: caldav.CompFilter{
			Name: "VCALENDAR",
			Comps: []caldav.CompFilter{{
				Name: "VEVENT",
				Props: []caldav.PropFilter{{
					Name:      ical.PropUID,
					TextMatch: &caldav.TextMatch{Text: uid},
				}},
			}},
		},
	}

	calendarURL := homeset + calendarName
	resp, err := client.QueryCalendar(ctx, calendarURL, query)
	if err != nil {
		return nil, fmt.Errorf("error getting calendar query: %v", err)
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("no events found with UID %s", uid)
	}
	return resp, nil
}

// tested
func ListTodos(ctx context.Context, client *caldav.Client, homeset, calendarName string) ([]caldav.CalendarObject, error) {
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
		return nil, fmt.Errorf("error getting calendar query: %v", err)
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("no todos found")
	}
	return resp, nil
}

// tested TODO split to smaller
func GetEvent(newEvent *Event) (*ical.Event, error) {
	event := ical.NewEvent()
	event.Name = newEvent.Name
	event.Props.SetText(ical.PropUID, newEvent.Uid)
	event.Props.SetDateTime(ical.PropDateTimeStamp, time.Now().UTC())
	event.Props.SetDateTime(ical.PropDateTimeStart, newEvent.DateTimeStart)
	event.Props.SetDateTime(ical.PropDateTimeEnd, newEvent.DateTimeEnd)
	// SetDTStart(event, newEvent)
	// SetDTEnd(event, newEvent)
	for _, attendee := range newEvent.Attendees {
		//AddAttendee(event, attendee)
		prop := ical.NewProp(ical.PropAttendee)
		prop.Params.Add(ical.ParamParticipationStatus, "NEEDS-ACTION")
		prop.Params.Add(ical.ParamRole, "REQ-PARTICIPANT")
		prop.Value = "mailto:" + attendee
		event.Props.Add(prop)
	}
	SetOrganizer(event, newEvent)
	AddAlarm(event, newEvent)
	return event, nil
}

func SetSummary(old *ical.Event, new *Event) error {
	// if empty, keep old value
	if new.Uid == "" {
		return nil
	}
	oldSummary, err := old.Props.Text(ical.PropSummary)
	if err != nil {
		return err
	}
	if oldSummary != new.Uid {
		old.Props.SetText(ical.PropUID, new.Uid)
	}
	return nil
}

// func SetDTStart(old *ical.Event, new *Event) {
// 	// if empty, keep old value
// 	if new.DateTimeStart.IsZero() {
// 		return
// 	}
// 	oldDTStart := old.Props.Get(ical.PropDateTimeStart).Value
// 	if oldDTStart != new.DateTimeStart.String() {
// 		old.Props.SetDateTime(ical.PropDateTimeStart, new.DateTimeStart)
// 	}
// }
// func SetDTEnd(old *ical.Event, new *Event) {
// 	// if empty, keep old value
// 	if new.DateTimeStart.IsZero() {
// 		return
// 	}
// 	oldDTEnd := old.Props.Get(ical.PropDateTimeEnd).Value
// 	if oldDTEnd != new.DateTimeEnd.String() {
// 		old.Props.SetDateTime(ical.PropDateTimeEnd, new.DateTimeEnd)
// 	}
// }

func AddAttendee(old *ical.Event, attendee string) {
	if attendee != "" {
		return
	}
	prop := ical.NewProp(ical.PropAttendee)
	prop.Params.Add(ical.ParamParticipationStatus, "NEEDS-ACTION")
	prop.Params.Add(ical.ParamRole, "REQ-PARTICIPANT")
	prop.Value = "mailto:" + attendee
	old.Props.Add(prop)
}

func SetOrganizer(old *ical.Event, new *Event) {
	if new.Attendees != nil {
		propOrg := ical.NewProp(ical.PropOrganizer)
		propOrg.Value = "mailto:" + new.Organizer
		old.Props.Add(propOrg)
	}
}

func AddAlarm(old *ical.Event, new *Event) {
	if new.Alarm != nil {
		alarm := ical.NewComponent(ical.CompAlarm)
		alarm.Props.SetText(ical.PropAction, new.Alarm.Action)
		alarm.Props.SetText(ical.PropTrigger, new.Alarm.Trigger)
		old.Children = append(old.Children, alarm)
	}

}

// tested TODO split to smaller
func GetTodo(newEvent *Event) (*ical.Event, error) {
	event := ical.NewEvent()
	event.Name = newEvent.Name
	event.Props.SetText(ical.PropUID, newEvent.Uid)
	err := SetSummary(event, newEvent)
	if err != nil {
		return nil, err
	}
	event.Props.SetDateTime(ical.PropDateTimeStamp, time.Now().UTC())
	//SetDTStart(event, newEvent)
	event.Props.SetDateTime(ical.PropDateTimeStart, newEvent.DateTimeStart)
	event.Props.SetDateTime(ical.PropDue, newEvent.DateTimeEnd)
	event.Props.SetText(ical.PropStatus, "NEEDS-ACTION")
	return event, nil
}

// tested TODO split to smaller
func GetRecurrentEvent(newRecEvent *ReccurentEvent) *ical.Event {
	event := ical.NewEvent()
	event.Name = newRecEvent.Event.Name
	event.Props.SetText(ical.PropUID, newRecEvent.Event.Uid)
	event.Props.SetDateTime(ical.PropDateTimeStamp, time.Now().UTC())
	event.Props.SetDateTime(ical.PropDateTimeStart, newRecEvent.Event.DateTimeStart)

	//SetDTStart(event, newRecEvent.Event)

	event.Props.SetRecurrenceRule(&rrule.ROption{
		Freq:       rrule.Frequency(newRecEvent.Frequency),
		Interval:   newRecEvent.Interval,
		Wkst:       rrule.MO,
		Count:      newRecEvent.Count,
		Until:      newRecEvent.Event.DateTimeEnd,
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

	for _, attendee := range newRecEvent.Event.Attendees {
		AddAttendee(event, attendee)
	}
	SetOrganizer(event, newRecEvent.Event)
	return event
}

// tested
func CreateEvent(ctx context.Context, client *caldav.Client, homeset string, calendarName string, event *ical.Event) error {
	calendar := ical.NewCalendar()
	calendar.Props.SetText(ical.PropVersion, "2.0")
	calendar.Props.SetText(ical.PropProductID, "-//trvita//EN")
	calendar.Props.SetText(ical.PropCalendarScale, "GREGORIAN")

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

func FindEventsWithExpand(ctx context.Context, httpClient webdav.HTTPClient, url, homeset, calendarName string, startTime, endTime time.Time) error {
	reqBody := fmt.Sprintf(`
	<?xml version="1.0" encoding="utf-8" ?>
	<C:calendar-query xmlns:D="DAV:"
					  xmlns:C="urn:ietf:params:xml:ns:caldav">
		<D:prop>
			<C:calendar-data>
				<C:expand start="%s" 
						  end="%s"/>
			</C:calendar-data>
		</D:prop>
		<C:filter>
			<C:comp-filter name="VCALENDAR">
				<C:comp-filter name="VEVENT">
					<C:time-range start="%s" 
								  end="%s"/>
				</C:comp-filter>
			</C:comp-filter>
		</C:filter>
	</C:calendar-query>`, startTime, endTime, startTime, endTime)

	calURL := url + homeset[8:] + calendarName
	req, err := http.NewRequest("REPORT", calURL, bytes.NewBufferString(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/xml; charset=\"utf-8\"")
	req.Header.Set("Depth", "1")

	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("REQUEST:\n%s\n", string(reqDump))
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMultiStatus {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("RESPONSE:\n%s\n", string(respDump))
	return nil
}

// tested
func Delete(ctx context.Context, client *caldav.Client, path string) error {
	err := client.RemoveAll(ctx, path)
	if err != nil {
		return err
	}
	return nil
}

func FindEvent(ctx context.Context, client *caldav.Client, eventURL, eventUID string) (*ical.Component, error) {
	obj, err := client.GetCalendarObject(ctx, eventURL)
	if err != nil {
		return nil, err
	}

	var foundComponent *ical.Component
	for _, comp := range obj.Data.Children {
		uid, err := comp.Props.Text(ical.PropUID)
		if err != nil {
			return nil, err
		}
		if uid == eventUID {
			foundComponent = comp
			break
		}
	}

	if foundComponent == nil {
		return nil, fmt.Errorf("event with UID %s not found", eventUID)
	}

	return foundComponent, nil
}

// tested
func ModifyAttendance(ctx context.Context, client *caldav.Client, homeset, calendarName, eventUID, eventPath string, mods *Modifications) error {
	eventURL := homeset[9:] + calendarName + "/" + eventPath + ".ics"
	comp, err := FindEvent(ctx, client, eventURL, eventUID)
	if err != nil {
		return err
	}
	newEventURL := homeset[9:] + mods.CalendarName + "/" + eventPath + ".ics"

	var att ical.Prop
	for _, att = range comp.Props.Values(ical.PropAttendee) {
		uri, err := att.URI()
		if err != nil {
			return err
		}
		if uri.String() == mods.Email {
			break
		}
	}

	if mods.PartStat == "DECLINED" {
		err = Delete(ctx, client, eventURL)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
	if mods.PartStat == "ACCEPTED" {
		att.Params.Set(ical.ParamParticipationStatus, mods.PartStat)
	}
	// add delegation support
	if !mods.LastModified.IsZero() {
		comp.Props.SetDateTime(ical.PropLastModified, mods.LastModified)
	}
	if mods.DelegateTo != "" {
		comp.Props.SetText(ical.ParamDelegatedTo, mods.DelegateTo)
	}
	comp.Props.SetDateTime(ical.PropLastModified, time.Now())

	calendar := ical.NewCalendar()
	calendar.Props.SetText(ical.PropVersion, "2.0")
	calendar.Props.SetText(ical.PropProductID, "-//trvita//EN")
	calendar.Props.SetText(ical.PropCalendarScale, "GREGORIAN")
	calendar.Children = append(calendar.Children, comp)
	_, err = client.PutCalendarObject(ctx, newEventURL, calendar)
	if err != nil {
		return err
	}
	return nil
}

// tested
func PutAttendee(ctx context.Context, client *caldav.Client, attendee, homeset, calendarName, eventUID, eventPath string) error {
	eventURL := homeset[9:] + calendarName + "/" + eventPath + ".ics"
	comp, err := FindEvent(ctx, client, eventURL, eventUID)
	if err != nil {
		return err
	}

	prop := ical.NewProp(ical.PropAttendee)
	prop.Params.Add(ical.ParamParticipationStatus, "NEEDS-ACTION")
	//prop.Params.Add(ical.ParamRole, "REQ-PARTICIPANT")
	prop.Value = "mailto:" + attendee
	comp.Props.Add(prop)

	calendar := ical.NewCalendar()
	calendar.Props.SetText(ical.PropVersion, "2.0")
	calendar.Props.SetText(ical.PropProductID, "-//trvita//EN")
	calendar.Props.SetText(ical.PropCalendarScale, "GREGORIAN")
	calendar.Children = append(calendar.Children, comp)

	_, err = client.PutCalendarObject(ctx, eventURL, calendar)
	if err != nil {
		return err
	}
	return nil
}

// TODO write func that requests attendees aka attendee syncronization
func LookUpAttendees() {

}

func UpdateEvent(ctx context.Context, client *caldav.Client, homeset, calendarName, eventUID, eventPath string) error {
	eventURL := homeset[9:] + "inbox/" + eventPath + ".ics"
	oldEventURL := homeset[9:] + calendarName + "/" + eventUID + ".ics"
	comp, err := FindEvent(ctx, client, eventURL, eventUID)
	if err != nil {
		return err
	}
	comp.Props.SetDateTime(ical.PropLastModified, time.Now())

	calendar := ical.NewCalendar()
	calendar.Props.SetText(ical.PropVersion, "2.0")
	calendar.Props.SetText(ical.PropProductID, "-//trvita//EN")
	calendar.Props.SetText(ical.PropCalendarScale, "GREGORIAN")
	calendar.Children = append(calendar.Children, comp)
	_, err = client.PutCalendarObject(ctx, oldEventURL, calendar)
	if err != nil {
		return err
	}
	return nil
}
