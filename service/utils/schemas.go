package utils

import (
	"errors"
	"regexp"
)

type Validable interface {
	Validate()
}

type User struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	DateOfBirth Date   `json:"dateOfBirth"`
}

func (u User) ValidateUsername() error {
	if l := len(u.Username); l < 3 || l > 16 {
		return errors.New("username lenght not valid")
	}
	m, e := regexp.MatchString("^[a-zA-Z0-9._]{3,16}$", u.Username)
	if !m {
		return e
	}
	return nil
}

func (u User) Validate() error {
	if e := u.ValidateUsername(); e != nil {
		return e
	}
	if l := len(u.Name); l < 3 || l > 25 {
		return errors.New("name lenght not valid")
	}
	if l := len(u.Surname); l < 3 || l > 25 {
		return errors.New("surname lenght not valid")
	}
	m, e := regexp.MatchString("^[a-zA-Z]{3,25}$", u.Name)
	if !m {
		return e
	}
	m, e = regexp.MatchString("^[a-zA-Z']{3,25}$", u.Surname)
	if !m {
		return e
	}
	e = u.DateOfBirth.validate()
	return e
}

type Ban struct {
	BannerID int `json:"bannerID"`
	BannedID int `json:"bannedID"`
}

type Follow struct {
	FollowerID int `json:"FollowerID"`
	FollowedID int `json:"FollowedID"`
}

type Like struct {
	UserID  int `json:"userID"`
	PhotoID int `json:"photoID"`
}
type Comment struct {
	Id      int    `json:"id"`
	UserID  int    `json:"userID"`
	PhotoID int    `json:"photoID"`
	Content string `json:"content"`
}

func (c Comment) Validate() error {
	if l := len(c.Content); l < 1 || l > 1024 {
		return errors.New("content lenght not valid")
	}
	return nil
}

type Date struct {
	Year  int    `json:"year"`
	Month string `json:"month"`
	Day   int    `json:"day"`
}

func (d Date) validate() error {
	validMonths := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	if d.Year < 1900 {
		return errors.New("bad year")
	}
	if !arrayContains(validMonths, d.Month) {
		return errors.New("bad month")
	}
	if d.Day < 1 || d.Day > 31 {
		return errors.New("bad day")
	}

	return nil
}

type Timestamp struct {
	Date    Date `json:"date"`
	Hour    int  `json:"hour"`
	Minutes int  `json:"minutes"`
	Seconds int  `json:"seconds"`
}

func (d Timestamp) Validate() error {
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
	UserId          int       `json:"userId"`
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
	if e := p.UploadTimestamp.Validate(); e != nil {
		return e
	}
	if l := len(p.Caption); l < 1 || l > 100 {
		return errors.New("caption lenght not valid")
	}
	return nil
}

type UserProfile struct {
	User      User    `json:"user"`
	Followed  []User  `json:"followed"`
	Following []User  `json:"following"`
	Stream    []Photo `json:"stream"`
}
