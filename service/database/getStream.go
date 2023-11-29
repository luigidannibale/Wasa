package database

import (
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

func (db *appdbimpl) GetStream(userID int) ([]utils.Photo, string, error) {

	var stream []utils.Photo
	// return stream,"Could't find the user", errors.New("UserNotFound")
	return stream, "An error occured on the server", errors.New("InternalServerError")
	return stream, "Stream successfully taken ", nil
}
