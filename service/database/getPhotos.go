package database

import (
	"database/sql"
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetPhotos(userID int) ([]utils.Photo, string, error) {
	var stream []utils.Photo
	rows, e := db.c.Query(`SELECT Id
						FROM Photos
						WHERE userID = ?`, userID)
	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return stream, "Couldn't find any photo", ErrNotFound
		}
		return stream, e.Error(), ErrInternalServerError
	}
	var photoIDs []int
	defer rows.Close()
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return stream, e.Error(), ErrInternalServerError
		}
		photoIDs = append(photoIDs, id)
	}

	for i := 0; i < len(photoIDs); i++ {
		p, _, e := db.GetPhoto(photoIDs[i])
		if e != nil {
			return stream, e.Error(), ErrInternalServerError
		}
		stream = append(stream, p)
	}

	return stream, "Photos found successfully", nil
}
