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
	var bannedIDs []int
	defer rows.Close()
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return banned, e.Error(), ErrInternalServerError
		}
		bannedIDs = append(bannedIDs, id)
	}

	for i := 0; i < len(bannedIDs); i++ {
		u, _, e := db.GetUser(bannedIDs[i])
		if e != nil {
			return banned, e.Error(), ErrInternalServerError
		}
		banned = append(banned, u)
	}

	return banned, "List of banned found successfully", nil
}
