package dtime_test

import (
	"fmt"
	"time"

	"github.com/rwxrob/bonzai/dtime"
)

const format = "2006-01-02 15:04 -0700"

func ExampleDayOfWeek() {
	then, _ := time.Parse("2006-01-02 15:04 -0700", "2020-05-13 14:34 -0500")
	fmt.Println(then)
	fmt.Println(dtime.DayOfWeek(&then, "tue"))
	fmt.Println(dtime.DayOfWeek(&then, "fri"))
	again, _ := time.Parse("2006-01-02 15:04 -0700", "2020-05-11 14:34 -0500")
	fmt.Println(again)
	fmt.Println(dtime.DayOfWeek(&again, "tue"))
	fmt.Println(dtime.DayOfWeek(&again, "fri"))
	// Output:
	// 2020-05-13 14:34:00 -0500 -0500
	// 2020-05-12 00:00:00 -0500 -0500
	// 2020-05-15 00:00:00 -0500 -0500
	// 2020-05-11 14:34:00 -0500 -0500
	// 2020-05-12 00:00:00 -0500 -0500
	// 2020-05-15 00:00:00 -0500 -0500
}

func ExampleMinuteOf() {
	t, _ := time.Parse("2006-01-02 15:04:05 -0700", "2020-05-13 14:34:56 -0500")
	fmt.Println(dtime.MinuteOf(&t))

	// Output:
	// 2020-05-13 14:34:00 -0500 -0500
}

func ExampleHourOf() {
	t, _ := time.Parse("2006-01-02 15:04:05 -0700", "2020-05-13 14:34:56 -0500")
	fmt.Println(dtime.HourOf(&t))

	// Output:
	// 2020-05-13 14:00:00 -0500 -0500
}

func ExampleDayOf() {
	t, _ := time.Parse("2006-01-02 15:04:05 -0700", "2020-05-13 14:34:56 -0500")
	fmt.Println(dtime.DayOf(&t))

	// Output:
	// 2020-05-13 00:00:00 -0500 -0500
}

func ExampleWeekOf() {
	t, _ := time.Parse("2006-01-02 15:04:05 -0700", "2020-05-13 14:34:56 -0500")
	fmt.Println(dtime.WeekOf(&t))

	// Output:
	// 2020-05-11 00:00:00 -0500 -0500
}

func ExampleMonthOf() {
	t, _ := time.Parse("2006-01-02 15:04:05 -0700", "2020-05-13 14:34:56 -0500")
	fmt.Println(dtime.MonthOf(&t))

	// Output:
	// 2020-05-01 00:00:00 -0500 -0500
}

func ExampleYearOf() {
	t, _ := time.Parse("2006-01-02 15:04:05 -0700", "2020-05-13 14:34:56 -0500")
	fmt.Println(dtime.YearOf(&t))

	// Output:
	// 2020-01-01 00:00:00 -0500 -0500
}

/*
func ExampleUntil() {

	now := time.Now().Add(1 * time.Hour)
	fmt.Println(dtime.Until(dtime.NextHourOf, &now))

	// Output:
	// ignored

}
*/
