package dtime

import (
	"fmt"
	"testing"
	"time"

	"github.com/rwxrob/bonzai/tinout"
)

const FMT = `2006-01-02 15:04:05 -0700`

var then, _ = time.Parse(FMT, "2020-05-13 14:34:54 -0500")

func TestSpan(t *testing.T) {
	spec, _ := tinout.Load("testdata/dtime.yaml")
	DefaultTime = &then

	result := spec.Check(func(tt *tinout.Test) bool {
		f, l := Span(tt.I)
		fs := "<nil>"
		ls := "<nil>"
		if f != nil {
			fs = f.Format(FMT)
		}
		if l != nil {
			ls = l.Format(FMT)
		}
		tt.Got = fmt.Sprintf("%v, %v", fs, ls)
		return tt.Passing()
	})
	if result != nil {
		t.Fatal(result.State())
	}
}

func TestMinuteOf(t *testing.T) {
	got := MinuteOf(&then)
	want, _ := time.Parse("2006-01-02 15:04:05 -0700", "2020-05-13 14:34:00 -0500")
	if want.Unix() != got.Unix() {
		t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
	}
}

func TestHourOf(t *testing.T) {
	got := HourOf(&then)
	want, _ := time.Parse("2006-01-02 15:04:05 -0700", "2020-05-13 14:00:00 -0500")
	if want.Unix() != got.Unix() {
		t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
	}
}

func TestDayOf(t *testing.T) {
	got := DayOf(&then)
	want, _ := time.Parse("2006-01-02 15:04:05 -0700", "2020-05-13 00:00:00 -0500")
	if want.Unix() != got.Unix() {
		t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
	}
}

func TestWeekOf(t *testing.T) {
	got := WeekOf(&then)
	want, _ := time.Parse("2006-01-02 15:04:05 -0700", "2020-05-11 00:00:00 -0500")
	if want.Unix() != got.Unix() {
		t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
	}
}

func TestMonthOf(t *testing.T) {
	got := MonthOf(&then)
	want, _ := time.Parse("2006-01-02 15:04:05 -0700", "2020-05-01 00:00:00 -0500")
	if want.Unix() != got.Unix() {
		t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
	}
}

func TestYearOf(t *testing.T) {
	got := YearOf(&then)
	want, _ := time.Parse("2006-01-02 15:04:05 -0700", "2020-01-01 00:00:00 -0500")
	if want.Unix() != got.Unix() {
		t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
	}
}

func TestToday(t *testing.T) {
	t.Log(Today())
}

func TestJanuaryOf(t *testing.T) {
	got := []*time.Time{
		JanuaryOf(&then),
		MonthOfYear(&then, "jan"),
		MonthOfYear(&then, "Jan"),
		MonthOfYear(&then, "january"),
		MonthOfYear(&then, "January"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-01-01 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestFebruaryOf(t *testing.T) {
	got := []*time.Time{
		FebruaryOf(&then),
		MonthOfYear(&then, "feb"),
		MonthOfYear(&then, "Feb"),
		MonthOfYear(&then, "february"),
		MonthOfYear(&then, "February"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-02-01 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestMarchOf(t *testing.T) {
	got := []*time.Time{
		MarchOf(&then),
		MonthOfYear(&then, "mar"),
		MonthOfYear(&then, "Mar"),
		MonthOfYear(&then, "march"),
		MonthOfYear(&then, "March"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-03-01 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestAprilOf(t *testing.T) {
	got := []*time.Time{
		AprilOf(&then),
		MonthOfYear(&then, "apr"),
		MonthOfYear(&then, "Apr"),
		MonthOfYear(&then, "april"),
		MonthOfYear(&then, "April"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-04-01 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestMayOf(t *testing.T) {
	got := []*time.Time{
		MayOf(&then),
		MonthOfYear(&then, "may"),
		MonthOfYear(&then, "May"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-05-01 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestJuneOf(t *testing.T) {
	got := []*time.Time{
		JuneOf(&then),
		MonthOfYear(&then, "jun"),
		MonthOfYear(&then, "Jun"),
		MonthOfYear(&then, "june"),
		MonthOfYear(&then, "June"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-06-01 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestJulyOf(t *testing.T) {
	got := []*time.Time{
		JulyOf(&then),
		MonthOfYear(&then, "jul"),
		MonthOfYear(&then, "Jul"),
		MonthOfYear(&then, "july"),
		MonthOfYear(&then, "July"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-07-01 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestAugustOf(t *testing.T) {
	got := []*time.Time{
		AugustOf(&then),
		MonthOfYear(&then, "aug"),
		MonthOfYear(&then, "Aug"),
		MonthOfYear(&then, "august"),
		MonthOfYear(&then, "August"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-08-01 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSeptemberOf(t *testing.T) {
	got := []*time.Time{
		SeptemberOf(&then),
		MonthOfYear(&then, "sep"),
		MonthOfYear(&then, "Sep"),
		MonthOfYear(&then, "september"),
		MonthOfYear(&then, "September"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-09-01 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestOctoberOf(t *testing.T) {
	got := []*time.Time{
		OctoberOf(&then),
		MonthOfYear(&then, "oct"),
		MonthOfYear(&then, "Oct"),
		MonthOfYear(&then, "october"),
		MonthOfYear(&then, "October"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-10-01 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestNovemberOf(t *testing.T) {
	got := []*time.Time{
		NovemberOf(&then),
		MonthOfYear(&then, "nov"),
		MonthOfYear(&then, "Nov"),
		MonthOfYear(&then, "november"),
		MonthOfYear(&then, "November"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-11-01 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestDecemberOf(t *testing.T) {
	got := []*time.Time{
		DecemberOf(&then),
		MonthOfYear(&then, "dec"),
		MonthOfYear(&then, "Dec"),
		MonthOfYear(&then, "december"),
		MonthOfYear(&then, "December"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-12-01 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestMonthOfYear(t *testing.T) {
	got := MonthOfYear(&then, "foo")
	if got != nil {
		t.Fail()
	}

}

func TestSameTimeInMonthOfYear(t *testing.T) {
	got := SameTimeInMonthOfYear(&then, "foo")
	if got != nil {
		t.Fail()
	}
}

func TestDayOfWeek(t *testing.T) {
	got := DayOfWeek(&then, "foo")
	if got != nil {
		t.Fail()
	}
}

func TestSameTimeOnDayOfWeek(t *testing.T) {
	got := SameTimeOnDayOfWeek(&then, "foo")
	if got != nil {
		t.Fail()
	}
}

func TestSameTimeInJanuaryOf(t *testing.T) {
	got := []*time.Time{
		SameTimeInJanuaryOf(&then),
		SameTimeInMonthOfYear(&then, "jan"),
		SameTimeInMonthOfYear(&then, "Jan"),
		SameTimeInMonthOfYear(&then, "january"),
		SameTimeInMonthOfYear(&then, "January"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-01-13 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeInFebruaryOf(t *testing.T) {
	got := []*time.Time{
		SameTimeInFebruaryOf(&then),
		SameTimeInMonthOfYear(&then, "feb"),
		SameTimeInMonthOfYear(&then, "Feb"),
		SameTimeInMonthOfYear(&then, "february"),
		SameTimeInMonthOfYear(&then, "February"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-02-13 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeInMarchOf(t *testing.T) {
	got := []*time.Time{
		SameTimeInMarchOf(&then),
		SameTimeInMonthOfYear(&then, "mar"),
		SameTimeInMonthOfYear(&then, "Mar"),
		SameTimeInMonthOfYear(&then, "march"),
		SameTimeInMonthOfYear(&then, "March"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-03-13 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeInAprilOf(t *testing.T) {
	got := []*time.Time{
		SameTimeInAprilOf(&then),
		SameTimeInMonthOfYear(&then, "apr"),
		SameTimeInMonthOfYear(&then, "Apr"),
		SameTimeInMonthOfYear(&then, "april"),
		SameTimeInMonthOfYear(&then, "April"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-04-13 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeInMayOf(t *testing.T) {
	got := []*time.Time{
		SameTimeInMayOf(&then),
		SameTimeInMonthOfYear(&then, "may"),
		SameTimeInMonthOfYear(&then, "May"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-05-13 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeInJuneOf(t *testing.T) {
	got := []*time.Time{
		SameTimeInJuneOf(&then),
		SameTimeInMonthOfYear(&then, "jun"),
		SameTimeInMonthOfYear(&then, "Jun"),
		SameTimeInMonthOfYear(&then, "june"),
		SameTimeInMonthOfYear(&then, "June"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-06-13 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeInJulyOf(t *testing.T) {
	got := []*time.Time{
		SameTimeInJulyOf(&then),
		SameTimeInMonthOfYear(&then, "jul"),
		SameTimeInMonthOfYear(&then, "Jul"),
		SameTimeInMonthOfYear(&then, "july"),
		SameTimeInMonthOfYear(&then, "July"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-07-13 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeInAugustOf(t *testing.T) {
	got := []*time.Time{
		SameTimeInAugustOf(&then),
		SameTimeInMonthOfYear(&then, "aug"),
		SameTimeInMonthOfYear(&then, "Aug"),
		SameTimeInMonthOfYear(&then, "august"),
		SameTimeInMonthOfYear(&then, "August"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-08-13 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeInSeptemberOf(t *testing.T) {
	got := []*time.Time{
		SameTimeInSeptemberOf(&then),
		SameTimeInMonthOfYear(&then, "sep"),
		SameTimeInMonthOfYear(&then, "Sep"),
		SameTimeInMonthOfYear(&then, "september"),
		SameTimeInMonthOfYear(&then, "September"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-09-13 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeInOctoberOf(t *testing.T) {
	got := []*time.Time{
		SameTimeInOctoberOf(&then),
		SameTimeInMonthOfYear(&then, "oct"),
		SameTimeInMonthOfYear(&then, "Oct"),
		SameTimeInMonthOfYear(&then, "october"),
		SameTimeInMonthOfYear(&then, "October"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-10-13 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeInNovemberOf(t *testing.T) {
	got := []*time.Time{
		SameTimeInNovemberOf(&then),
		SameTimeInMonthOfYear(&then, "nov"),
		SameTimeInMonthOfYear(&then, "Nov"),
		SameTimeInMonthOfYear(&then, "november"),
		SameTimeInMonthOfYear(&then, "November"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-11-13 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeInDecemberOf(t *testing.T) {
	got := []*time.Time{
		SameTimeInDecemberOf(&then),
		SameTimeInMonthOfYear(&then, "dec"),
		SameTimeInMonthOfYear(&then, "Dec"),
		SameTimeInMonthOfYear(&then, "december"),
		SameTimeInMonthOfYear(&then, "December"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-12-13 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

/*
      May 2020
Su Mo Tu We Th Fr Sa
                1  2
 3  4  5  6  7  8  9
10 11 12 13 14 15 16
17 18 19 20 21 22 23
24 25 26 27 28 29 30
31
*/

func TestMondayOf(t *testing.T) {
	got := []*time.Time{
		MondayOf(&then),
		DayOfWeek(&then, "mon"),
		DayOfWeek(&then, "Mon"),
		DayOfWeek(&then, "monday"),
		DayOfWeek(&then, "Monday"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-05-11 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestTuesdayOf(t *testing.T) {
	got := []*time.Time{
		TuesdayOf(&then),
		DayOfWeek(&then, "tue"),
		DayOfWeek(&then, "Tue"),
		DayOfWeek(&then, "tuesday"),
		DayOfWeek(&then, "Tuesday"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-05-12 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestWednesdayOf(t *testing.T) {
	got := []*time.Time{
		WednesdayOf(&then),
		DayOfWeek(&then, "wed"),
		DayOfWeek(&then, "Wed"),
		DayOfWeek(&then, "wednesday"),
		DayOfWeek(&then, "Wednesday"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-05-13 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestThursdayOf(t *testing.T) {
	got := []*time.Time{
		ThursdayOf(&then),
		DayOfWeek(&then, "thu"),
		DayOfWeek(&then, "Thu"),
		DayOfWeek(&then, "thursday"),
		DayOfWeek(&then, "Thursday"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-05-14 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestFridayOf(t *testing.T) {
	got := []*time.Time{
		FridayOf(&then),
		DayOfWeek(&then, "fri"),
		DayOfWeek(&then, "Fri"),
		DayOfWeek(&then, "friday"),
		DayOfWeek(&then, "Friday"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-05-15 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSaturdayOf(t *testing.T) {
	got := []*time.Time{
		SaturdayOf(&then),
		DayOfWeek(&then, "sat"),
		DayOfWeek(&then, "Sat"),
		DayOfWeek(&then, "saturday"),
		DayOfWeek(&then, "Saturday"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-05-16 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSundayOf(t *testing.T) {
	got := []*time.Time{
		SundayOf(&then),
		DayOfWeek(&then, "sun"),
		DayOfWeek(&then, "Sun"),
		DayOfWeek(&then, "sunday"),
		DayOfWeek(&then, "Sunday"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-05-17 00:00:00 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeMondayOf(t *testing.T) {
	got := []*time.Time{
		SameTimeOnMondayOf(&then),
		SameTimeOnDayOfWeek(&then, "mon"),
		SameTimeOnDayOfWeek(&then, "Mon"),
		SameTimeOnDayOfWeek(&then, "monday"),
		SameTimeOnDayOfWeek(&then, "Monday"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-05-11 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeTuesdayOf(t *testing.T) {
	got := []*time.Time{
		SameTimeOnTuesdayOf(&then),
		SameTimeOnDayOfWeek(&then, "tue"),
		SameTimeOnDayOfWeek(&then, "Tue"),
		SameTimeOnDayOfWeek(&then, "tuesday"),
		SameTimeOnDayOfWeek(&then, "Tuesday"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-05-12 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeWednesdayOf(t *testing.T) {
	got := []*time.Time{
		SameTimeOnWednesdayOf(&then),
		SameTimeOnDayOfWeek(&then, "wed"),
		SameTimeOnDayOfWeek(&then, "Wed"),
		SameTimeOnDayOfWeek(&then, "wednesday"),
		SameTimeOnDayOfWeek(&then, "Wednesday"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-05-13 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeThursdayOf(t *testing.T) {
	got := []*time.Time{
		SameTimeOnThursdayOf(&then),
		SameTimeOnDayOfWeek(&then, "thu"),
		SameTimeOnDayOfWeek(&then, "Thu"),
		SameTimeOnDayOfWeek(&then, "thursday"),
		SameTimeOnDayOfWeek(&then, "Thursday"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-05-14 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeFridayOf(t *testing.T) {
	got := []*time.Time{
		SameTimeOnFridayOf(&then),
		SameTimeOnDayOfWeek(&then, "fri"),
		SameTimeOnDayOfWeek(&then, "Fri"),
		SameTimeOnDayOfWeek(&then, "friday"),
		SameTimeOnDayOfWeek(&then, "Friday"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-05-15 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeSaturdayOf(t *testing.T) {
	got := []*time.Time{
		SameTimeOnSaturdayOf(&then),
		SameTimeOnDayOfWeek(&then, "sat"),
		SameTimeOnDayOfWeek(&then, "Sat"),
		SameTimeOnDayOfWeek(&then, "saturday"),
		SameTimeOnDayOfWeek(&then, "Saturday"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-05-16 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}

func TestSameTimeSundayOf(t *testing.T) {
	got := []*time.Time{
		SameTimeOnSundayOf(&then),
		SameTimeOnDayOfWeek(&then, "sun"),
		SameTimeOnDayOfWeek(&then, "Sun"),
		SameTimeOnDayOfWeek(&then, "sunday"),
		SameTimeOnDayOfWeek(&then, "Sunday"),
	}
	for _, g := range got {
		want, _ := time.Parse(FMT, "2020-05-17 14:34:54 -0500")
		if want.Unix() != g.Unix() {
			t.Fatalf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}
