package database

import (
	"database/sql"
	"errors"
	"fmt"

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
	if rows.Err() != nil {
		return banned, rows.Err().Error(), ErrInternalServerError
	}
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
			return banned, err.Error(), ErrInternalServerError
		}
		banned = append(banned, u)
	}
	return banned, "List of banned found successfully", nil
}

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) DeleteBan(ban utils.Ban) (string, error) {
	userID, userToUnbanID := ban.BannerID, ban.BannedID

	res, err := db.c.Exec(`DELETE FROM Bans WHERE BannerID = ? AND BannedID = ?`, userID, userToUnbanID)

	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return fmt.Sprintf("User %d was not banned by %d", userToUnbanID, userID), ErrNotFound
	}
	if err != nil {
		return err.Error(), ErrInternalServerError
	}

	return fmt.Sprintf("User %d has unbanned %d", userID, userToUnbanID), nil

}

/*
Errors that can be returned: (InternalServerError, AlreadyDone)
*/
func (db *appdbimpl) CreateBan(ban utils.Ban) (string, error) {
	userID, userToBanID := ban.BannerID, ban.BannedID
	var retS string
	var retE error = nil
	tx, er := db.c.Begin()
	if er != nil {
		return er.Error(), ErrInternalServerError
	}
	res, err := tx.Exec(`INSERT OR IGNORE INTO Bans(BannerID,BannedID)
									VALUES (?,?)`, userID, userToBanID)

	if x, y := res.RowsAffected(); x == 0 && y == nil {
		retS, retE = "Already banned", ErrAlreadyDone
	}
	if err != nil {
		retS, retE = err.Error(), ErrInternalServerError
	}
	res, err = tx.Exec(`DELETE FROM Follows 
			WHERE (FollowerID = ? AND FollowedID = ?)
			OR	(FollowerID = ? AND FollowedID = ?)`,
		userID, userToBanID, userToBanID, userID)

	if err != nil {
		retS, retE = err.Error(), ErrInternalServerError
	}
	if retE != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return rbErr.Error(), ErrInternalServerError
		}
		return retS, retE
	}

	e := tx.Commit()
	if e != nil {
		return e.Error(), ErrInternalServerError
	}

	return fmt.Sprintf("User %d banned %d", userID, userToBanID), nil

}

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) CheckBan(banner int, banned int) error {
	var b utils.Ban
	b.BannerID = banner
	b.BannedID = banned
	_, er := db.GetBan(b)
	return er
}

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
