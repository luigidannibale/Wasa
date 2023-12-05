package database

import (
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

func (db *appdbimpl) CreateComment(comment utils.Comment) (int, string, error) {
	var commentID int

	_, s, e := db.GetUser(comment.UserID)
	if e != nil {
		if e.Error() == "NotFound" {
			return commentID, "Couldn't find the user", errors.New("UserNotFound")
		}
		if e.Error() == "InternalServerError" {
			return commentID, s, errors.New("InternalServerError")
		}
	}
	/*
		err := db.c.QueryRow(`INSERT OR IGNORE INTO Photos(UserID,Image,Caption,UploadTimestamp)
								VALUES (?,?,?,?)
								RETURNING Id`, photo.UserId, photo.Image, photo.Caption, utils.TimestampToString(photo.UploadTimestamp)).Scan(&photoID)
		if err != nil {
			return photoID, "An error occcurred on the server" + err.Error(), errors.New("InternalServerError")
		}*/
	return commentID, "Photo uploaded successfully", nil
}
