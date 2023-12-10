package utils

import (
	"strconv"
	"strings"
)

func arrayContains[T comparable](array []T, val T) bool {
	for _, v := range array {
		// checking if the array contains the given value
		if v == val {
			// changing the boolean variable
			return true
		}
	}
	return false
}

func remove[T interface{}](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}
func StringToDate(s string) Date {
	// s is in "dd-month-yyyy" format
	var d Date
	x := strings.Split(s, "-")
	day, _ := strconv.Atoi(x[0])
	year, _ := strconv.Atoi(x[2])
	d.Day = day
	d.Month = x[1]
	d.Year = year
	return d
}

func DateToString(d Date) string {
	var s string = "dd-month-yyyy"
	s = strconv.Itoa(d.Day) + "-" + d.Month + "-" + strconv.Itoa(d.Year)
	return s
}
func StringToTimestamp(s string) Timestamp {
	// s is in "dd-month-yyyy-hh-mm-ss" format
	var t Timestamp
	x := strings.Split(s, "-")
	day, _ := strconv.Atoi(x[0])
	year, _ := strconv.Atoi(x[2])
	hour, _ := strconv.Atoi(x[3])
	minutes, _ := strconv.Atoi(x[4])
	seconds, _ := strconv.Atoi(x[5])
	t.Date.Day = day
	t.Date.Month = x[1]
	t.Date.Year = year
	t.Hour = hour
	t.Minutes = minutes
	t.Seconds = seconds
	return t
}

func TimestampToString(t Timestamp) string {
	var s string = "dd-month-yyyy-hh-mm-ss"
	s = DateToString(t.Date) + "-" + strconv.Itoa(t.Hour) + "-" + strconv.Itoa(t.Minutes) + "-" + strconv.Itoa(t.Seconds)
	return s
}
func Now() Timestamp {
	var t Timestamp
	return t
}
