// Package dtime enables easy human input of times, dates, and durations. It
// also includes many convenience functions for rouding time duration
// boundaries as is frequently needed for scheduling and time-based search
// applications.
//
// Pointers to time.Time are used throughout the package since <nil> is usually
// a more desirable zero value than that for time.Time without a pointer.  */
package dtime

import (
	"fmt"
	"time"
)

const (
	SECOND float64 = 1000000000
	MINUTE float64 = 60000000000
	HOUR   float64 = 3600000000000
	DAY    float64 = 86400000000000
	WEEK   float64 = 604800000000000
	YEAR   float64 = 31536000000000000
)

// DefaultTime can be set instead of time.Now().
var DefaultTime *time.Time

func _deftime() *time.Time {
	if DefaultTime == nil {
		now := time.Now()
		return &now
	}
	return DefaultTime
}

func dump(i interface{}) {
	fmt.Printf("%v\n", i)
}

// Span parses the htime format string and returns one or two time pointers.
// The first is the start time, the second is the last time. If a date and time
// are detected in the string the first is set. If the offset is detected in
// the string the second will be set. If no date and time are detected
// DefaultTime or time.Now() is assumed.
//
// Spans the main function in this package and provides a minimal format for
// entering date and time information in a practical way. It's primary use case
// is when a user needs to enter such data quickly and regularly from the
// command line or into mobile and other devices where human input speed is
// limited to what can be tapped out on the screen. The characters used in the
// formatting only characters that appear on most default keyboards without
// shifting or switching to symbolic input. The package provides no method of
// specifying timezone, which falls out of the scope of this package. All times
// are therefore assumed to be local. For a precise specification of the format
// see the htime.peg file included with the package source code.
func Span(s string) (first *time.Time, last *time.Time) {
	p := new(spanParser)
	p.Buffer = s
	p.Init()
	p.Parse()
	// p.Execute() // no need with -inline
	if !p.start.IsZero() {
		first = &p.start
	}
	if p.offset != 0 {
		t := p.start.Add(time.Duration(int64(p.offset)))
		last = &t
	}
	if !p.stop.IsZero() {
		last = &p.stop
	}
	return
}

// MinuteOf returns the start of the given minute.
func MinuteOf(t *time.Time) *time.Time {
	nt := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	return &nt
}

// HourOf returns the start of the given hour.
func HourOf(t *time.Time) *time.Time {
	nt := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	return &nt
}

// NextHourOf returns the start of the given hour.
func NextHourOf(t *time.Time) *time.Time {
	nt := time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+1, 0, 0, 0, t.Location())
	return &nt
}

// DayOf returns the start of the given day.
func DayOf(t *time.Time) *time.Time {
	nt := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return &nt
}

// WeekOf returns the start of the given week.
func WeekOf(t *time.Time) *time.Time {
	return MondayOf(t)
}

// MonthOf returns the start of the month.
func MonthOf(t *time.Time) *time.Time {
	nt := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	return &nt
}

// YearOf returns the start of the month.
func YearOf(t *time.Time) *time.Time {
	nt := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	return &nt
}

// Today returns the start of the current date.
func Today() *time.Time {
	return DayOf(_deftime())
}

// Tomorrow returns the start of the next day.
func Tomorrow() *time.Time {
	t := *_deftime()
	t = t.Add(time.Hour * 24)
	return DayOf(&t)
}

// Yesterday returns the start of the previous day.
func Yesterday() *time.Time {
	t := *_deftime()
	t = t.Add(time.Hour * -24)
	return DayOf(&t)
}

// JanuaryOf returns the beginning of the first day of January of the given
// year.
func JanuaryOf(t *time.Time) *time.Time {
	return MonthOf(SameTimeInJanuaryOf(t))
}

// FebruaryOf returns the beginning of the first day of February of the given
// year.
func FebruaryOf(t *time.Time) *time.Time {
	return MonthOf(SameTimeInFebruaryOf(t))
}

// MarchOf returns the beginning of the first day of March of the given
// year.
func MarchOf(t *time.Time) *time.Time {
	return MonthOf(SameTimeInMarchOf(t))
}

// AprilOf returns the beginning of the first day of April of the given
// year.
func AprilOf(t *time.Time) *time.Time {
	return MonthOf(SameTimeInAprilOf(t))
}

// MayOf returns the beginning of the first day of May of the given
// year.
func MayOf(t *time.Time) *time.Time {
	return MonthOf(SameTimeInMayOf(t))
}

// JuneOf returns the beginning of the first day of June of the given
// year.
func JuneOf(t *time.Time) *time.Time {
	return MonthOf(SameTimeInJuneOf(t))
}

// JulyOf returns the beginning of the first day of July of the given
// year.
func JulyOf(t *time.Time) *time.Time {
	return MonthOf(SameTimeInJulyOf(t))
}

// AugustOf returns the beginning of the first day of August of the given
// year.
func AugustOf(t *time.Time) *time.Time {
	return MonthOf(SameTimeInAugustOf(t))
}

// SeptemberOf returns the beginning of the first day of September of the given
// year.
func SeptemberOf(t *time.Time) *time.Time {
	return MonthOf(SameTimeInSeptemberOf(t))
}

// OctoberOf returns the beginning of the first day of October of the given
// year.
func OctoberOf(t *time.Time) *time.Time {
	return MonthOf(SameTimeInOctoberOf(t))
}

// NovemberOf returns the beginning of the first day of November of the given
// year.
func NovemberOf(t *time.Time) *time.Time {
	return MonthOf(SameTimeInNovemberOf(t))
}

// DecemberOf returns the beginning of the first day of December of the given
// year.
func DecemberOf(t *time.Time) *time.Time {
	return MonthOf(SameTimeInDecemberOf(t))
}

func samemonth(t *time.Time, month int) *time.Time {
	d := time.Date(
		t.Year(),
		time.Month(month),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond(),
		t.Location(),
	)
	return &d
}

// SameTimeInJanuaryOf returns the exact same month day and time but for the
// month of January instead.
func SameTimeInJanuaryOf(t *time.Time) *time.Time {
	return samemonth(t, 1)
}

// SameTimeInFebruaryOf returns the exact same month day and time but for the
// month of February instead.
func SameTimeInFebruaryOf(t *time.Time) *time.Time {
	return samemonth(t, 2)
}

// SameTimeInMarchOf returns the exact same month day and time but for the
// month of March instead.
func SameTimeInMarchOf(t *time.Time) *time.Time {
	return samemonth(t, 3)
}

// SameTimeInAprilOf returns the exact same month day and time but for the
// month of April instead.
func SameTimeInAprilOf(t *time.Time) *time.Time {
	return samemonth(t, 4)
}

// SameTimeInMayOf returns the exact same month day and time but for the
// month of May instead.
func SameTimeInMayOf(t *time.Time) *time.Time {
	return samemonth(t, 5)
}

// SameTimeInJuneOf returns the exact same month day and time but for the
// month of June instead.
func SameTimeInJuneOf(t *time.Time) *time.Time {
	return samemonth(t, 6)
}

// SameTimeInJulyOf returns the exact same month day and time but for the
// month of July instead.
func SameTimeInJulyOf(t *time.Time) *time.Time {
	return samemonth(t, 7)
}

// SameTimeInAugustOf returns the exact same month day and time but for the
// month of August instead.
func SameTimeInAugustOf(t *time.Time) *time.Time {
	return samemonth(t, 8)
}

// SameTimeInSeptemberOf returns the exact same month day and time but for the
// month of September instead.
func SameTimeInSeptemberOf(t *time.Time) *time.Time {
	return samemonth(t, 9)
}

// SameTimeInOctoberOf returns the exact same month day and time but for the
// month of October instead.
func SameTimeInOctoberOf(t *time.Time) *time.Time {
	return samemonth(t, 10)
}

// SameTimeInNovemberOf returns the exact same month day and time but for the
// month of November instead.
func SameTimeInNovemberOf(t *time.Time) *time.Time {
	return samemonth(t, 11)
}

// SameTimeInDecemberOf returns the exact same month day and time but for the
// month of December instead.
func SameTimeInDecemberOf(t *time.Time) *time.Time {
	return samemonth(t, 12)
}

// MonthOfYear returns the beginning of the first day of the specified month of the
// given year.
func MonthOfYear(t *time.Time, month string) *time.Time {
	switch month {
	case "jan", "january", "Jan", "January":
		return JanuaryOf(t)
	case "feb", "Feb", "february", "February":
		return FebruaryOf(t)
	case "mar", "Mar", "march", "March":
		return MarchOf(t)
	case "apr", "Apr", "april", "April":
		return AprilOf(t)
	case "may", "May":
		return MayOf(t)
	case "jun", "Jun", "june", "June":
		return JuneOf(t)
	case "jul", "Jul", "july", "July":
		return JulyOf(t)
	case "aug", "Aug", "august", "August":
		return AugustOf(t)
	case "sep", "Sep", "september", "September":
		return SeptemberOf(t)
	case "oct", "Oct", "october", "October":
		return OctoberOf(t)
	case "nov", "Nov", "november", "November":
		return NovemberOf(t)
	case "dec", "Dec", "december", "December":
		return DecemberOf(t)
	}
	return nil
}

// SameTimeInMonthOfYear returns the exact same month day and time but for the
// specified month instead.
func SameTimeInMonthOfYear(t *time.Time, month string) *time.Time {
	switch month {
	case "jan", "Jan", "january", "January":
		return SameTimeInJanuaryOf(t)
	case "feb", "Feb", "february", "February":
		return SameTimeInFebruaryOf(t)
	case "mar", "Mar", "march", "March":
		return SameTimeInMarchOf(t)
	case "apr", "Apr", "april", "April":
		return SameTimeInAprilOf(t)
	case "may", "May":
		return SameTimeInMayOf(t)
	case "jun", "Jun", "june", "June":
		return SameTimeInJuneOf(t)
	case "jul", "Jul", "july", "July":
		return SameTimeInJulyOf(t)
	case "aug", "Aug", "august", "August":
		return SameTimeInAugustOf(t)
	case "sep", "Sep", "september", "September":
		return SameTimeInSeptemberOf(t)
	case "oct", "Oct", "october", "October":
		return SameTimeInOctoberOf(t)
	case "nov", "Nov", "november", "November":
		return SameTimeInNovemberOf(t)
	case "dec", "Dec", "december", "December":
		return SameTimeInDecemberOf(t)
	}
	return nil
}

// MondayOf returns the Monday of the week passed as time rounded to the
// beginning of the day.
func MondayOf(t *time.Time) *time.Time {
	return DayOf(SameTimeOnMondayOf(t))
}

// TuesdayOf returns the Tuesday of the week passed as time rounded to the
// beginning of the day.
func TuesdayOf(t *time.Time) *time.Time {
	return DayOf(SameTimeOnTuesdayOf(t))
}

// WednesdayOf returns the Wednesday of the week passed as time rounded to the
// beginning of the day.
func WednesdayOf(t *time.Time) *time.Time {
	return DayOf(SameTimeOnWednesdayOf(t))
}

// ThursdayOf returns the Thursday of the week passed as time rounded to the
// beginning of the day.
func ThursdayOf(t *time.Time) *time.Time {
	return DayOf(SameTimeOnThursdayOf(t))
}

// FridayOf returns the Friday of the week passed as time rounded to the
// beginning of the day.
func FridayOf(t *time.Time) *time.Time {
	return DayOf(SameTimeOnFridayOf(t))
}

// SaturdayOf returns the Saturday of the week passed as time rounded to the
// beginning of the day.
func SaturdayOf(t *time.Time) *time.Time {
	return DayOf(SameTimeOnSaturdayOf(t))
}

// SundayOf returns the Sunday of the week passed as time rounded to the
// beginning of the day.
func SundayOf(t *time.Time) *time.Time {
	return DayOf(SameTimeOnSundayOf(t))
}

// SameTimeOnMondayOf returns the exact same time but on the Monday of the week
// indicated. For beginning of day use without SameTimeOn.
func SameTimeOnMondayOf(t *time.Time) *time.Time {
	nt := t.Add((time.Duration(time.Monday-t.Weekday()) * 24 * time.Hour))
	return &nt
}

// SameTimeOnTuesdayOf returns the exact same time but on the Tuesday of the week
// indicated. For beginning of day use without SameTimeOn.
func SameTimeOnTuesdayOf(t *time.Time) *time.Time {
	nt := t.Add((time.Duration(time.Tuesday-t.Weekday()) * 24 * time.Hour))
	return &nt
}

// SameTimeOnWednesdayOf returns the exact same time but on the Wednesday of the week
// indicated. For beginning of day use without SameTimeOn.
func SameTimeOnWednesdayOf(t *time.Time) *time.Time {
	nt := t.Add((time.Duration(time.Wednesday-t.Weekday()) * 24 * time.Hour))
	return &nt
}

// SameTimeOnThursdayOf returns the exact same time but on the Thursday of the week
// indicated. For beginning of day use without SameTimeOn.
func SameTimeOnThursdayOf(t *time.Time) *time.Time {
	nt := t.Add((time.Duration(time.Thursday-t.Weekday()) * 24 * time.Hour))
	return &nt
}

// SameTimeOnFridayOf returns the exact same time but on the Friday of the week
// indicated. For beginning of day use without SameTimeOn.
func SameTimeOnFridayOf(t *time.Time) *time.Time {
	nt := t.Add((time.Duration(time.Friday-t.Weekday()) * 24 * time.Hour))
	return &nt
}

// SameTimeOnSaturdayOf returns the exact same time but on the Saturday of the week
// indicated. For beginning of day use without SameTimeOn.
func SameTimeOnSaturdayOf(t *time.Time) *time.Time {
	nt := t.Add((time.Duration(time.Saturday-t.Weekday()) * 24 * time.Hour))
	return &nt
}

// SameTimeOnSundayOf returns the exact same time but on the Sunday of the week
// indicated. For beginning of day use without SameTimeOn.
func SameTimeOnSundayOf(t *time.Time) *time.Time {
	nt := t.Add((time.Duration(time.Sunday-t.Weekday()+7) * 24 * time.Hour))
	return &nt
}

// DayOfWeek returns the day of the week passed rounded to the beginning of the
// weekday indicated.
func DayOfWeek(t *time.Time, day string) *time.Time {
	switch day {
	case "mon", "monday", "Mon", "Monday":
		return MondayOf(t)
	case "tue", "tuesday", "Tue", "Tuesday":
		return TuesdayOf(t)
	case "wed", "wednesday", "Wed", "Wednesday":
		return WednesdayOf(t)
	case "thu", "thursday", "Thu", "Thursday":
		return ThursdayOf(t)
	case "fri", "friday", "Fri", "Friday":
		return FridayOf(t)
	case "sat", "saturday", "Sat", "Saturday":
		return SaturdayOf(t)
	case "sun", "sunday", "Sun", "Sunday":
		return SundayOf(t)
	}
	return nil
}

// SameTimeOnDayOfWeek returns the day of the week passed rounded to the
// beginning of the week day indicated.
func SameTimeOnDayOfWeek(t *time.Time, day string) *time.Time {
	switch day {
	case "mon", "monday", "Mon", "Monday":
		return SameTimeOnMondayOf(t)
	case "tue", "tuesday", "Tue", "Tuesday":
		return SameTimeOnTuesdayOf(t)
	case "wed", "wednesday", "Wed", "Wednesday":
		return SameTimeOnWednesdayOf(t)
	case "thu", "thursday", "Thu", "Thursday":
		return SameTimeOnThursdayOf(t)
	case "fri", "friday", "Fri", "Friday":
		return SameTimeOnFridayOf(t)
	case "sat", "saturday", "Sat", "Saturday":
		return SameTimeOnSaturdayOf(t)
	case "sun", "sunday", "Sun", "Sunday":
		return SameTimeOnSundayOf(t)
	}
	return nil
}

// Until takes a time and a function that takes a time and returns a time and
// converts the result into a duration offset between them. Returns nil
// if unable to determine. This is useful when a duration is needed for
// any of the functions in this package that match that signature. The
// duration is positive or negative depending on the relation between
// the times.
func Until(fn func(*time.Time) *time.Time, t *time.Time) time.Duration {
	return (*(fn(t))).Sub(*t)
}
