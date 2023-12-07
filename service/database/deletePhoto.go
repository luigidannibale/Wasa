package database

import (
	"fmt"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) DeletePhoto(photo utils.Photo) (string, error) {

	res, err := db.c.Exec(`DELETE FROM Photos WHERE Id = ?`, photo.Id)

	if err != nil {
		return err.Error(), ErrInternalServerError
	}
	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return fmt.Sprintf("Can't find a photo with id %d", photo.Id), ErrNotFound
	}
	return fmt.Sprintf("Photo with id %d has been deleted", photo.Id), nil
}
