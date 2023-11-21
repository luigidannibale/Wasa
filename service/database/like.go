package database

import (
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

// Takes photoID and puts the like
func (db *appdbimpl) Like(photoID int, like utils.Like) (string, error) {

	// return "The photo was already liked", errors.New("AlreadyLiked")
	// return "Could't find the user", errors.New("UserNotFound")
	// return "Could't find the photo", errors.New("PhotoNotFound")
	return "An error occured on the server", errors.New("InternalServerError")
	//return "Photo liked successfully", nil
}
