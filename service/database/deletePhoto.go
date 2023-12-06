package database

import (
	"fmt"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) DeletePhoto(photo utils.Photo) (string, error) {

	res, err := db.c.Exec(`DELETE FROM Photos WHERE PhotoID = ?`, photo.Id)

	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return fmt.Sprintf("Can't find a photo with id %d", photo.Id), ErrNotFound
	}
	if err != nil {
		return err.Error(), ErrInternalServerError
	}
	return fmt.Sprintf("Photo with id %d has been deleted", photo.Id), nil
}
