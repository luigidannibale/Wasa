package database

//import "errors"

func (db *appdbimpl) Unfollow(userID int, userToUnfollowID int) (string, error) {

	// return "Could't find the user", errors.New("UserNotFound")
	// return "Could't find the user to unfollow", errors.New("FollowedNotFound")
	// return "An error occured on the server", errors.New("InternalServerError")
	return "User unfollowed successfully", nil
}
