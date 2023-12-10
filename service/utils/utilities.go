package utils

import (
	"sort"
)

type By func(p1, p2 *Photo) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(photos []Photo) {
	ps := &photoSorter{
		photos: photos,
		by:     by,
	}
	sort.Sort(ps)
}

type photoSorter struct {
	photos []Photo
	by     func(p1, p2 *Photo) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *photoSorter) Len() int {
	return len(s.photos)
}

// Swap is part of sort.Interface.
func (s *photoSorter) Swap(i, j int) {
	s.photos[i], s.photos[j] = s.photos[j], s.photos[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *photoSorter) Less(i, j int) bool {
	return s.by(&s.photos[i], &s.photos[j])
}

func SortStreamByLast(unsortedStream []Photo) ([]Photo, error) {
	ts := func(p1, p2 *Photo) bool {
		return p1.UploadTimestamp.After(p2.UploadTimestamp)
	}
	By(ts).Sort(unsortedStream)
	return unsortedStream, nil
}

/*
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
}*/
