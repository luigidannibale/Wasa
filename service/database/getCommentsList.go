package database

import (
	"database/sql"
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetCommentsList(photoID int) ([]utils.Comment, string, error) {
	var comments []utils.Comment
	rows, e := db.c.Query(`SELECT Id, userID, Content
						FROM Comments
						WHERE PhotoID = ?`, photoID)
	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return comments, "Couldn't find any comment", ErrNotFound
		}
		return comments, e.Error(), ErrInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		var c utils.Comment
		err := rows.Scan(&c.Id, &c.UserID, &c.Content)
		if err != nil {
			return comments, e.Error(), ErrInternalServerError
		}
		c.PhotoID = photoID
		comments = append(comments, c)
	}

	return comments, "List of comments found successfully", nil
}
