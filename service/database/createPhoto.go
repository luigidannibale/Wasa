package database

import (
	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (InternalServerError)
*/
func (db *appdbimpl) CreatePhoto(photo utils.Photo) (int, string, error) {
	var photoID int

	err := db.c.QueryRow(`INSERT INTO Photos(UserID,Image,Caption,UploadTimestamp)
							VALUES (?,?,?,?)
							RETURNING Id`, photo.UserId, photo.Image, photo.Caption, utils.TimestampToString(photo.UploadTimestamp)).Scan(&photoID)
	if err != nil {
		return photoID, err.Error(), InternalServerError
	}
	return photoID, "Photo uploaded successfully", nil
}
