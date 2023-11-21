package utils

import (
	"errors"
)

type Validable interface {
	Validate()
}

type User struct {
	id       int
	Username string
}

func (u User) Validate() error {
	if l := len(u.Username); l < 3 || l > 16 {
		return errors.New("Username lenght not valid")
	}

	return nil
}

type Photo struct {
	id    int
	Image string
}

type Like struct {
	UserID int
}
type Comment struct {
	id      int
	UserID  int
	Content string
}
type Date struct {
	Year    int
	Month   string
	Day     int
	Hour    int
	Minutes int
	Seconds int
}

func (d Date) validate() error {
	validMonths := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

	if !arrayContains(validMonths, d.Month) {
		return errors.New("Bad month")
	}
	if d.Day < 1 || d.Day > 31 {
		return errors.New("Bad day")
	}
	if d.Hour < 0 || d.Hour > 23 {
		return errors.New("Bad hour")
	}
	if d.Minutes < 0 || d.Minutes > 59 {
		return errors.New("Bad minutes")
	}
	if d.Seconds < 0 || d.Seconds > 59 {
		return errors.New("Bad seconds")
	}

	return nil
}

type UserProfile struct {
	user      User
	followed  []User
	following []User
	stream    []Photo
}
