package database

import (
	"database/sql"
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetBannedList(userID int) ([]utils.User, string, error) {
	var banned []utils.User
	rows, e := db.c.Query(`SELECT BannedID
						FROM Bans
						WHERE BannerID = ?`, userID)
	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return banned, "Couldn't find any banned", ErrNotFound
		}
		return banned, e.Error(), ErrInternalServerError
	}
	e = rows.Scan(&banned)
	if e != nil {
		return banned, e.Error(), ErrInternalServerError
	}
	return banned, "List of banned found successfully", nil
}
