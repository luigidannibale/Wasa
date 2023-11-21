package database

//import "errors"

//Takes userID and photoID and removess the user like to the photo
func (db *appdbimpl) Unlike(userID int, photoID int) (string, error) {

	// return "Could't find the user", errors.New("UserNotFound")
	// return "Could't find the photo", errors.New("PhotoNotFound")
	// return "An error occured on the server", errors.New("InternalServerError")
	return "Photo unliked successfully", nil
}
