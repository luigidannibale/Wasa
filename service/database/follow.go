package database

//import "errors"

func (db *appdbimpl) Follow(userID int, userToFollowID int) (string, error) {

	// return "Could't find the user", errors.New("UserNotFound")
	// return "Could't find the user to follow", errors.New("FollowedNotFound")
	// return "An error occured on the server", errors.New("InternalServerError")
	// return "The user was already followed", errors.New("AlreadyFollowed")
	return "User followed successfully", nil
}
