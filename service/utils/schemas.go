package utils

import (
	"errors"
)

type Validable interface {
	Validate()
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
}

// to insert regex validation
func (u User) Validate() error {
	if l := len(u.Username); l < 3 || l > 16 {
		return errors.New("username lenght not valid")
	}
	if l := len(u.Name); l < 3 || l > 25 {
		return errors.New("name lenght not valid")
	}
	if l := len(u.Surname); l < 3 || l > 25 {
		return errors.New("surname lenght not valid")
	}
	return nil
}

type Like struct {
	Id int `json:"id"`
}
type Comment struct {
	Id      int    `json:"id"`
	UserID  int    `json:"userID"`
	Content string `json:"content"`
}
type Timestamp struct {
	Year    int    `json:"year"`
	Month   string `json:"month"`
	Day     int    `json:"day"`
	Hour    int    `json:"hour"`
	Minutes int    `json:"minutes"`
	Seconds int    `json:"seconds"`
}

func (d Timestamp) validate() error {
	validMonths := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

	if !arrayContains(validMonths, d.Month) {
		return errors.New("bad month")
	}
	if d.Day < 1 || d.Day > 31 {
		return errors.New("bad day")
	}
	if d.Hour < 0 || d.Hour > 23 {
		return errors.New("bad hour")
	}
	if d.Minutes < 0 || d.Minutes > 59 {
		return errors.New("bad minutes")
	}
	if d.Seconds < 0 || d.Seconds > 59 {
		return errors.New("bad seconds")
	}

	return nil
}

type Photo struct {
	Id              int       `json:"id"`
	Image           string    `json:"image"`
	Caption         string    `json:"caption"`
	UploadTimestamp Timestamp `json:"uploadTimestamp"`
}

func (p Photo) Validate() error {
	if l := len(p.Image); l < 0 || l > 40960 {
		return errors.New("photo not acceptable")
	}
	if l := len(p.Caption); l < 0 || l > 100 {
		return errors.New("caption not acceptable")
	}
	if e := p.UploadTimestamp.validate(); e != nil {
		return e
	}
	return nil
}

type UserProfile struct {
	user      User
	followed  []User
	following []User
	stream    []Photo
}
