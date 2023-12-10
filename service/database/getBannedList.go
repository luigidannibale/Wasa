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
	rows, e := db.c.Query(`SELECT Id, Username, Name, Surname, DateOfBirth
						FROM Bans
						JOIN Users ON BannedID = Id
						WHERE BannerID = ? `, userID)
	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return banned, "Couldn't find any banned", ErrNotFound
		}
		return banned, e.Error(), ErrInternalServerError
	}

	defer rows.Close()
	for rows.Next() {
		var u utils.User
		err := rows.Scan(&u.Id, &u.Username, &u.Name, &u.Surname, &u.DateOfBirth)
		if err != nil {
			return banned, e.Error(), ErrInternalServerError
		}
		banned = append(banned, u)
	}
	return banned, "List of banned found successfully", nil
}
