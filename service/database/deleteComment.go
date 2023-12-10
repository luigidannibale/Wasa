package database

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
