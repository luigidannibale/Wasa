package database

import (
	"database/sql"
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetBan(ban utils.Ban) (string, error) {

	e := db.c.QueryRow(`SELECT BannerID, BannedID
						FROM Bans
						WHERE BannerID = ? AND BannedID = ?`, ban.BannerID, ban.BannedID)
	if e.Err() == nil {
		return "Ban found successfully", nil
	}
	if errors.Is(e.Err(), sql.ErrNoRows) {
		return "Couldn't find the ban", ErrNotFound
	}
	return e.Err().Error(), ErrInternalServerError
}
