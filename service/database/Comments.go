package database

import (
	"database/sql"
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (InternalServerError)
*/
func (db *appdbimpl) CreateComment(comment utils.Comment) (int, string, error) {
	var commentID int

	err := db.c.QueryRow(`INSERT INTO Comments(UserID,PhotoID,Content) VALUES (?,?,?) RETURNING Id`, comment.UserID, comment.PhotoID, comment.Content).Scan(&commentID)
	if err != nil {
		return commentID, err.Error(), ErrInternalServerError
	}
	return commentID, "Comment posted successfully", nil
}

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) DeleteComment(commentID int) (string, error) {

	res, err := db.c.Exec(`DELETE 
						FROM Comments
						WHERE Id = ?`, commentID)

	if err != nil {
		return err.Error(), ErrInternalServerError
	}
	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return "No commment with such id", ErrNotFound
	}
	return "Comment deleted successfully", nil
}

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetComment(commentID int) (utils.Comment, string, error) {
	var comment utils.Comment
	e := db.c.QueryRow(`SELECT Id,UserID,PhotoID,Content
						FROM Comments
						WHERE Id = ?`, commentID).Scan(&comment.Id, &comment.UserID, &comment.PhotoID, &comment.Content)
	if e == nil {
		return comment, "Comment found successfully", nil
	}
	if errors.Is(e, sql.ErrNoRows) {
		return comment, "Couldn't find the user", ErrNotFound
	}
	return comment, e.Error(), ErrInternalServerError
}

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
