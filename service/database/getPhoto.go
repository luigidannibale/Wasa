package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetPhoto(photoID int) (utils.Photo, string, error) {
	var photo utils.Photo
	var caption sql.NullString
	var ts string
	e := db.c.QueryRow(`SELECT Id,UserID,Image,Caption,UploadTimestamp
						FROM Photos
						WHERE Id = ?`, photoID).Scan(&photo.Id, &photo.UserId, &photo.Image, &caption, &ts)
	if caption.Valid {
		photo.Caption = caption.String
	}
	photo.UploadTimestamp = utils.StringToTimestamp(ts)
	if e == nil {
		return photo, "Photo found successfully", nil
	}
	fmt.Println("err" + e.Error())
	if errors.Is(e, sql.ErrNoRows) {
		return photo, "Couldn't find the photo", ErrNotFound
	}
	return photo, e.Error(), ErrInternalServerError
}
