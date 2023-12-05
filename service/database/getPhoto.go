package database

import (
	"database/sql"
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Returns user searched by the id,
a string message,
and an error
*/
func (db *appdbimpl) GetPhoto(photoID int) (utils.Photo, string, error) {
	var photo utils.Photo
	var caption sql.NullString
	e := db.c.QueryRow(`SELECT Id,UserID,Image,Caption,UploadTimestamp
						FROM Photos
						WHERE Id = ?`, photoID).Scan(&photo.Id, &photo.UserId, &photo.Image, &caption, &photo.UploadTimestamp)
	if caption.Valid {
		photo.Caption = caption.String
	}
	if e == nil {
		return photo, "Photo found successfully", nil
	}
	if errors.Is(e, sql.ErrNoRows) {
		return photo, "Couldn't find the user", errors.New("NotFound")
	}
	return photo, "An error occured on the server : " + e.Error(), errors.New("InternalServerError")
}
