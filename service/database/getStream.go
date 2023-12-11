package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetStream(userID int) ([]utils.Photo, string, error) {
	var stream []utils.Photo
	rows, e := db.c.Query(`SELECT Id, Image, Caption, UploadTimestamp
						FROM Photos
						JOIN Follows ON userID = FollowedID 
						WHERE FollowerID = ?
						ORDER BY UploadTimestamp DESC`, userID)
	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return stream, "Couldn't find any photo", ErrNotFound
		}
		return stream, e.Error(), ErrInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		var p utils.Photo
		var ts string
		err := rows.Scan(&p.Id, &p.Image, &p.Caption, &ts)
		if err != nil {
			return stream, e.Error(), ErrInternalServerError
		}
		ut, er := time.Parse(time.Layout, ts)
		if er != nil {
			return stream, er.Error(), ErrInternalServerError
		}
		p.UploadTimestamp = ut
		p.UserId = userID
		stream = append(stream, p)
	}
	return stream, "Photos found successfully", nil
}
