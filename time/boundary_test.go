package time

import (
	"fmt"
	"time"
)

var testCST, _ = time.LoadLocation(AsiaShanghai)

func ExampleStartOfCentury() {
	fmt.Println(StartOfCentury(time.Date(2023, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(StartOfCentury(time.Date(2023, 5, 5, 23, 59, 57, 123, testCST)).String())
	// Output:
	// 2000-01-01 00:00:00 +0000 UTC
	// 2000-01-01 00:00:00 +0800 CST
}

func ExampleEndOfCentury() {
	fmt.Println(EndOfCentury(time.Date(2023, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(EndOfCentury(time.Date(2023, 5, 5, 23, 59, 57, 123, testCST)).String())
	// Output:
	// 2099-12-31 23:59:59.999999999 +0000 UTC
	// 2099-12-31 23:59:59.999999999 +0800 CST
}

func ExampleStartOfDecade() {
	fmt.Println(StartOfDecade(time.Date(2023, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(StartOfDecade(time.Date(2023, 5, 5, 23, 59, 57, 123, testCST)).String())
	// Output:
	// 2020-01-01 00:00:00 +0000 UTC
	// 2020-01-01 00:00:00 +0800 CST
}

func ExampleEndOfDecade() {
	fmt.Println(EndOfDecade(time.Date(2023, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(EndOfDecade(time.Date(2023, 5, 5, 23, 59, 57, 123, testCST)).String())
	// Output:
	// 2029-12-31 23:59:59.999999999 +0000 UTC
	// 2029-12-31 23:59:59.999999999 +0800 CST
}

func ExampleStartOfYear() {
	fmt.Println(StartOfYear(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(StartOfYear(time.Date(2020, 5, 5, 23, 59, 57, 123, testCST)).String())
	// Output:
	// 2020-01-01 00:00:00 +0000 UTC
	// 2020-01-01 00:00:00 +0800 CST
}

func ExampleEndOfYear() {
	fmt.Println(EndOfYear(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(EndOfYear(time.Date(2020, 5, 5, 23, 59, 57, 123, testCST)).String())
	// Output:
	// 2020-12-31 23:59:59.999999999 +0000 UTC
	// 2020-12-31 23:59:59.999999999 +0800 CST
}

func ExampleStartOfQuarter() {
	fmt.Println(StartOfQuarter(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(StartOfQuarter(time.Date(2020, 5, 5, 23, 59, 57, 123, testCST)).String())
	fmt.Println(StartOfQuarter(time.Date(2020, 9, 10, 14, 36, 34, 123, time.UTC)).String())
	fmt.Println(StartOfQuarter(time.Date(2020, 11, 17, 8, 7, 6, 123, testCST)).String())
	// Output:
	// 2020-01-01 00:00:00 +0000 UTC
	// 2020-04-01 00:00:00 +0800 CST
	// 2020-07-01 00:00:00 +0000 UTC
	// 2020-10-01 00:00:00 +0800 CST
}

func ExampleEndOfQuarter() {
	fmt.Println(EndOfQuarter(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(EndOfQuarter(time.Date(2020, 5, 5, 23, 59, 57, 123, testCST)).String())
	fmt.Println(EndOfQuarter(time.Date(2020, 9, 10, 14, 36, 34, 123, time.UTC)).String())
	fmt.Println(EndOfQuarter(time.Date(2020, 11, 17, 8, 7, 6, 123, testCST)).String())
	// Output:
	// 2020-03-31 23:59:59.999999999 +0000 UTC
	// 2020-06-30 23:59:59.999999999 +0800 CST
	// 2020-09-30 23:59:59.999999999 +0000 UTC
	// 2020-12-31 23:59:59.999999999 +0800 CST
}

func ExampleStartOfMonth() {
	fmt.Println(StartOfMonth(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(StartOfMonth(time.Date(2020, 2, 14, 23, 59, 57, 123, testCST)).String())
	// Output:
	// 2020-02-01 00:00:00 +0000 UTC
	// 2020-02-01 00:00:00 +0800 CST
}

func ExampleEndOfMonth() {
	fmt.Println(EndOfMonth(time.Date(2019, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(EndOfMonth(time.Date(2019, 2, 14, 23, 59, 57, 123, testCST)).String())
	fmt.Println(EndOfMonth(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(EndOfMonth(time.Date(2020, 2, 14, 23, 59, 57, 123, testCST)).String())
	// Output:
	// 2019-02-28 23:59:59.999999999 +0000 UTC
	// 2019-02-28 23:59:59.999999999 +0800 CST
	// 2020-02-29 23:59:59.999999999 +0000 UTC
	// 2020-02-29 23:59:59.999999999 +0800 CST
}

func ExampleStartOfWeek() {
	fmt.Println(StartOfWeek(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC), time.Sunday).String())
	fmt.Println(StartOfWeek(time.Date(2020, 2, 14, 23, 59, 57, 123, testCST), time.Sunday).String())
	fmt.Println(StartOfWeek(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC), time.Monday).String())
	fmt.Println(StartOfWeek(time.Date(2020, 2, 14, 23, 59, 57, 123, testCST), time.Monday).String())
	// Output:
	// 2020-02-02 00:00:00 +0000 UTC
	// 2020-02-09 00:00:00 +0800 CST
	// 2020-01-27 00:00:00 +0000 UTC
	// 2020-02-10 00:00:00 +0800 CST
}

func ExampleEndOfWeek() {
	fmt.Println(EndOfWeek(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC), time.Sunday).String())
	fmt.Println(EndOfWeek(time.Date(2020, 2, 14, 23, 59, 57, 123, testCST), time.Sunday).String())
	fmt.Println(EndOfWeek(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC), time.Monday).String())
	fmt.Println(EndOfWeek(time.Date(2020, 2, 14, 23, 59, 57, 123, testCST), time.Monday).String())
	// Output:
	// 2020-02-08 23:59:59.999999999 +0000 UTC
	// 2020-02-15 23:59:59.999999999 +0800 CST
	// 2020-02-02 23:59:59.999999999 +0000 UTC
	// 2020-02-16 23:59:59.999999999 +0800 CST
}

func ExampleStartOfDay() {
	fmt.Println(StartOfDay(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(StartOfDay(time.Date(2020, 2, 14, 23, 59, 57, 123, testCST)).String())
	// Output:
	// 2020-02-02 00:00:00 +0000 UTC
	// 2020-02-14 00:00:00 +0800 CST
}

func ExampleEndOfDay() {
	fmt.Println(EndOfDay(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(EndOfDay(time.Date(2020, 2, 14, 23, 59, 57, 123, testCST)).String())
	// Output:
	// 2020-02-02 23:59:59.999999999 +0000 UTC
	// 2020-02-14 23:59:59.999999999 +0800 CST
}

func ExampleStartOfHour() {
	fmt.Println(StartOfHour(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(StartOfHour(time.Date(2020, 2, 14, 23, 59, 57, 123, testCST)).String())
	// Output:
	// 2020-02-02 12:00:00 +0000 UTC
	// 2020-02-14 23:00:00 +0800 CST
}

func ExampleEndOfHour() {
	fmt.Println(EndOfHour(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(EndOfHour(time.Date(2020, 2, 14, 23, 59, 57, 123, testCST)).String())
	// Output:
	// 2020-02-02 12:59:59.999999999 +0000 UTC
	// 2020-02-14 23:59:59.999999999 +0800 CST
}

func ExampleStartOfMinute() {
	fmt.Println(StartOfMinute(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(StartOfMinute(time.Date(2020, 2, 14, 23, 59, 57, 123, testCST)).String())
	// Output:
	// 2020-02-02 12:30:00 +0000 UTC
	// 2020-02-14 23:59:00 +0800 CST
}

func ExampleEndOfMinute() {
	fmt.Println(EndOfMinute(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(EndOfMinute(time.Date(2020, 2, 14, 23, 59, 57, 123, testCST)).String())
	// Output:
	// 2020-02-02 12:30:59.999999999 +0000 UTC
	// 2020-02-14 23:59:59.999999999 +0800 CST
}

func ExampleStartOfSecond() {
	fmt.Println(StartOfSecond(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(StartOfSecond(time.Date(2020, 2, 14, 23, 59, 57, 123, testCST)).String())
	// Output:
	// 2020-02-02 12:30:14 +0000 UTC
	// 2020-02-14 23:59:57 +0800 CST
}

func ExampleEndOfSecond() {
	fmt.Println(EndOfSecond(time.Date(2020, 2, 2, 12, 30, 14, 123, time.UTC)).String())
	fmt.Println(EndOfSecond(time.Date(2020, 2, 14, 23, 59, 57, 123, testCST)).String())
	// Output:
	// 2020-02-02 12:30:14.999999999 +0000 UTC
	// 2020-02-14 23:59:57.999999999 +0800 CST
}
