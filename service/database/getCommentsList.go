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
	rows, e := db.c.Query(`SELECT Id
						FROM Comments
						WHERE PhotoID = ?`, photoID)
	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return comments, "Couldn't find any comment", ErrNotFound
		}
		return comments, e.Error(), ErrInternalServerError
	}
	var commentIDs []int
	defer rows.Close()
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return comments, e.Error(), ErrInternalServerError
		}
		commentIDs = append(commentIDs, id)
	}

	for i := 0; i < len(commentIDs); i++ {
		u, _, e := db.GetComment(commentIDs[i])
		if e != nil {
			return comments, e.Error(), ErrInternalServerError
		}
		comments = append(comments, u)
	}

	return comments, "List of comments found successfully", nil
}
