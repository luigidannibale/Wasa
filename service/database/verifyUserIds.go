package database

import (
	"errors"
)

func (db *appdbimpl) VerifyUserIds(userIds []int) (int, error) {
	for i := 0; i < len(userIds); i++ {
		_, _, e := db.GetUser(userIds[i])
		if e != nil {
			if e.Error() == "NotFound" {
				return i, errors.New("NotFound")
			}
			if e.Error() == "InternalServerError" {
				return i, errors.New("InternalServerError")
			}
		}
	}
	return 0, nil

}
