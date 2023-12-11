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
	var b utils.Ban

	er := db.c.QueryRow(`SELECT BannerID
						FROM Bans
						WHERE BannerID = ? AND BannedID = ?`, ban.BannerID, ban.BannedID).Scan(&b.BannerID)

	if er == nil {
		return "Ban found successfully", nil
	}
	if errors.Is(er, sql.ErrNoRows) {
		return "Couldn't find the ban", ErrNotFound
	}
	return er.Error(), ErrInternalServerError
}
