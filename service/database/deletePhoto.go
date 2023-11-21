package database

func (db *appdbimpl) DeletePhoto(photoID int) (string, error) {
	// return "Could't find the photo", errors.New("PhotoNotFound")
	// return "An error occured on the server", errors.New("InternalServerError")
	return "Photo deleted successfully", nil
}
