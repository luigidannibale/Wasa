package database

import (
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

// Takes photoID and puts the like
func (db *appdbimpl) UploadPhoto(userID int, photo utils.Photo) (int, string, error) {

	var photoID int
	// return photoID,"Could't find the user", errors.New("UserNotFound")
	return photoID, "An error occured on the server", errors.New("InternalServerError")
	return photoID, "Photo posted successfully", nil
}
