package database

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) DeleteComment(commentID int) (string, error) {

	res, err := db.c.Exec(`DELETE OR IGNORE 
							FROM Comments
							WHERE CommentId = ?`, commentID)

	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return "No commment with such id", ErrNotFound
	}
	if err != nil {
		return err.Error(), ErrInternalServerError
	}
	return "Comment deleted successfully", nil
}
