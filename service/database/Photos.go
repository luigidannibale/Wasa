package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (InternalServerError)
*/
func (db *appdbimpl) CreatePhoto(photo utils.Photo) (int, string, error) {
	var photoID int

	err := db.c.QueryRow(`INSERT INTO Photos(UserID,Image,Caption,UploadTimestamp)
							VALUES (?,?,?,?)
							RETURNING Id`, photo.UserId, photo.Image, photo.Caption, photo.UploadTimestamp.Format(time.Layout)).Scan(&photoID)
	if err != nil {
		return photoID, err.Error(), ErrInternalServerError
	}
	return photoID, "Photo uploaded successfully", nil
}

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

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetPhoto(photoID int) (utils.Photo, string, error) {
	var photo utils.Photo
	var caption sql.NullString
	var ts string
	e := db.c.QueryRow(`SELECT Id,UserID,Image,Caption,UploadTimestamp
						FROM Photos
						WHERE Id = ?`, photoID).Scan(&photo.Id, &photo.UserId, &photo.Image, &caption, &ts)
	if caption.Valid {
		photo.Caption = caption.String
	}
	ut, er := time.Parse(time.Layout, ts)
	if er != nil {
		return photo, er.Error(), ErrInternalServerError
	}
	photo.UploadTimestamp = ut
	if e == nil {
		return photo, "Photo found successfully", nil
	}
	if errors.Is(e, sql.ErrNoRows) {
		return photo, "Couldn't find the photo", ErrNotFound
	}
	return photo, e.Error(), ErrInternalServerError
}

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

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetMyPhotos(userID int) ([]utils.Photo, string, error) {
	var stream []utils.Photo
	rows, e := db.c.Query(`SELECT Photos.Id, Image, Caption, UploadTimestamp
						FROM Photos		
						Where UserId = ?
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
