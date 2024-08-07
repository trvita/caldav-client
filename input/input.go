package input

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/trvita/caldav-client/mycal"
	"github.com/trvita/go-ical"
)

func String(r io.Reader, message string) (string, error) {
	reader := bufio.NewReader(r)
	if r == r {
		fmt.Print(message)
	}
	str, err := reader.ReadString('\n')
	str = strings.Trim(str, "\n")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(str), nil
}

func Commands(comms string) []string {
	commands := strings.Split(comms, "\n")
	return commands
}

func Int(r io.Reader, message string) (int, error) {
	var num int
	reader := bufio.NewReader(r)
	if r == r {
		fmt.Print(message)
	}
	_, err := fmt.Fscanf(reader, "%d\n", &num)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func Ints(r io.Reader, message string) ([]int, error) {
	reader := bufio.NewReader(r)
	if r == r {
		fmt.Print(message)
	}
	str, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	str = strings.TrimSpace(str)
	if str == "" {
		return nil, err
	}
	numbersStr := strings.Split(str, ",")
	var numbersInt []int

	for _, numStr := range numbersStr {
		numStr = strings.TrimSpace(numStr)
		numInt, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, err
		}
		numbersInt = append(numbersInt, numInt)
	}

	return numbersInt, nil
}

func Event(r io.Reader) (*mycal.Event, error) {
	var attendees []string
	var summary, organizer, name, startDate, startTime, endDate, endTime, action, trigger string
	var startDateTime, endDateTime time.Time

	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	summary, err = String(r, "Enter event summary: ")
	if err != nil {
		return nil, err
	}
	for {
		name, err = String(r, "Enter event type [event, todo]: ")
		if err != nil {
			return nil, err
		}
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
		startDate, err = String(r, "Enter event start date (YYYY.MM.DD): ")
		if err != nil {
			return nil, err
		}
		startTime, err = String(r, "Enter event start time (HH.MM.SS): ")
		if err != nil {
			return nil, err
		}
		startDateTime, err = time.Parse("2006.01.02 15.04.05", startDate+" "+startTime)
		if err != nil {
			fmt.Println("invalid start date/time format")
			continue
		}
		break
	}
	for {
		endDate, err = String(r, "Enter event end date (YYYY.MM.DD): ")
		if err != nil {
			return nil, err
		}
		endTime, err = String(r, "Enter event end time (HH.MM.SS): ")
		if err != nil {
			return nil, err
		}
		endDateTime, err = time.Parse("2006.01.02 15.04.05", endDate+" "+endTime)
		if err != nil {
			fmt.Println("invalid end date/time format")
			continue
		}
		break
	}
	for {
		attendee, err := String(r, "Enter attendee email (or 0 to finish): ")
		if err != nil {
			return nil, err
		}
		if attendee == "0" {
			break
		}
		attendees = append(attendees, attendee)
	}
	if attendees != nil {
		organizer, err = String(r, "Enter organizer email: ")
		if err != nil {
			return nil, err
		}

	}
	hasalarm, err := String(r, "Add alarm [y/n]: ")
	if err != nil {
		return nil, err
	}
	if hasalarm != "y" {
		return &mycal.Event{
			Name:          name,
			Uid:           uid.String(),
			Summary:       summary,
			DateTimeStart: startDateTime,
			DateTimeEnd:   endDateTime,
			Attendees:     attendees,
			Organizer:     organizer,
			Alarm:         nil,
		}, nil
	}
	if hasalarm == "y" {
		action, err := String(r, "Enter action: [d - display, e - email]: ")
		if err != nil {
			return nil, err
		}
		switch action {
		case "d":
			action = ical.ParamDisplay
		case "e":
			action = ical.ParamEmail
		}

	}
	return &mycal.Event{
		Name:          name,
		Uid:           uid.String(),
		Summary:       summary,
		DateTimeStart: startDateTime,
		DateTimeEnd:   endDateTime,
		Attendees:     attendees,
		Organizer:     organizer,
		Alarm: &mycal.Alarm{
			Action:  action,
			Trigger: trigger,
		},
	}, nil
}

func RecurrentEvent(r io.Reader) (*mycal.ReccurentEvent, error) {
	var attendees []string
	var byDay, byMonthDay, byYearDay, byMonth, byWeekNo, bySetPos, byHour []int
	var summary, name, startDate, startTime, freq, untilDate, untilTime, organizer string
	var startDateTime, untilDateTime time.Time
	var frequency, interval, count, ans int
	name = "VEVENT"
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	summary, err = String(r, "Enter event summary: ")
	if err != nil {
		return nil, err
	}
	for {
		startDate, err = String(r, "Enter event start date (YYYY.MM.DD): ")
		if err != nil {
			return nil, err
		}
		startTime, err = String(r, "Enter event start time (HH.MM.SS): ")
		if err != nil {
			return nil, err
		}

		startDateTime, err = time.Parse("2006.01.02 15.04.05", startDate+" "+startTime)
		if err != nil {
			fmt.Println("invalid start date/time format")
			continue
		}
		break
	}
	cont := true
	for cont {
		freq, err = String(r, "Enter frequency [Y, MO, W, D, H, MI, S]: ")
		if err != nil {
			return nil, err
		}
		switch strings.ToUpper(freq) {
		case "Y":
			frequency = 0
			cont = false
		case "MO":
			frequency = 1
			cont = false
		case "W":
			frequency = 2
			cont = false
		case "D":
			frequency = 3
			cont = false
		case "H":
			frequency = 4
			cont = false
		case "MI":
			frequency = 5
			cont = false
		case "S":
			frequency = 6
			cont = false
		}
	}
	interval, err = Int(r, "Enter interval: ")
	if err != nil {
		return nil, err
	}
	ans, err = Int(r, "Count, until or skip? [1/2/0]: ")
	if err != nil {
		return nil, err
	}
	switch ans {
	case 1:
		count, err = Int(r, "Enter count: ")
		if err != nil {
			return nil, err
		}
	case 2:
		for {
			untilDate, err = String(r, "Enter event start date (YYYY.MM.DD): ")
			if err != nil {
				return nil, err
			}
			untilTime, err = String(r, "Enter event start time (HH.MM.SS): ")
			if err != nil {
				return nil, err
			}

			untilDateTime, err = time.Parse("2006.01.02 15.04.05", untilDate+" "+untilTime)
			if err != nil {
				fmt.Println("invalid start date/time format")
				continue
			}
			break
		}
	}

	byDay, err = Ints(r, "Enter by days [num of day in week, num of day in year]: ")
	if err != nil {
		return nil, err
	}
	byMonthDay, err = Ints(r, "Enter by month days: ")
	if err != nil {
		return nil, err
	}
	byYearDay, err = Ints(r, "Enter by year days: ")
	if err != nil {
		return nil, err
	}
	byMonth, err = Ints(r, "Enter by months: ")
	if err != nil {
		return nil, err
	}
	byWeekNo, err = Ints(r, "Enter by week numbers: ")
	if err != nil {
		return nil, err
	}
	bySetPos, err = Ints(r, "Enter position by set: ")
	if err != nil {
		return nil, err
	}
	byHour, err = Ints(r, "Enter by hour numbers")
	if err != nil {
		return nil, err
	}

	for {
		attendee, err := String(r, "Enter attendee email (or 0 to finish): ")
		if err != nil {
			return nil, err
		}
		if attendee == "0" {
			break
		}
		attendees = append(attendees, attendee)
	}
	if attendees != nil {
		organizer, err = String(r, "Enter organizer email: ")
		if err != nil {
			return nil, err
		}

	}
	return &mycal.ReccurentEvent{
		Event: &mycal.Event{
			Name:          name,
			Summary:       summary,
			Uid:           uid.String(),
			DateTimeStart: startDateTime,
			DateTimeEnd:   untilDateTime,
			Attendees:     attendees,
			Organizer:     organizer,
		},
		Frequency:  frequency,
		Count:      count,
		Interval:   interval,
		ByDay:      byDay,
		ByMonthDay: byMonthDay,
		ByYearDay:  byYearDay,
		ByMonth:    byMonth,
		ByWeekNo:   byWeekNo,
		BySetPos:   bySetPos,
		ByHour:     byHour}, nil
}

func Modifications(r io.Reader) (*mycal.Modifications, error) {
	var partstat, delegateto, calendarName, email, answer string
	var err error

	email, err = String(r, "Enter your email: ")
	if err != nil {
		return nil, err
	}
	fmt.Println("-1")
	answer, err = String(r, "Accept, decline, delegate event? [y, n, d]: ")
	if err != nil {
		return nil, err
	}
	fmt.Println("0")
	switch answer {
	case "y":
		fmt.Println("1")
		partstat = "ACCEPTED"
		calendarName, err = String(r, "Enter which calendar event goes to: ")
		if err != nil {
			return nil, err
		}
	case "d":
		partstat = string(ical.ParamDelegatedTo)
		delegateto, err = String(r, "Enter who to delegate: ")
		if err != nil {
			return nil, err
		}
	default:
		fmt.Println("2")
		partstat = "DECLINED"
	}
	fmt.Println("3")
	return &mycal.Modifications{
		Email:        email,
		PartStat:     partstat,
		DelegateTo:   "mailto:" + delegateto,
		CalendarName: calendarName,
		LastModified: time.Now(),
	}, nil
}
