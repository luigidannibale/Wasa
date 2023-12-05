package database

import (
	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (InternalServerError)
*/
func (db *appdbimpl) CreateComment(comment utils.Comment) (int, string, error) {
	var commentID int

	err := db.c.QueryRow(`INSERT INTO Comments(UserID,PhotoID,Content) VALUES (?,?,?) RETURNING Id`, comment.UserID, comment.PhotoID, comment.Content).Scan(&commentID)
	if err != nil {
		return commentID, err.Error(), InternalServerError
	}
	return commentID, "Comment posted successfully", nil
}
