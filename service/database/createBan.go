package database

import (
	"fmt"

	"github.com/luigidannibale/Wasa/service/utils"
)

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
		tx.Rollback()
		return retS, retE
	}

	e := tx.Commit()
	if e != nil {
		return e.Error(), ErrInternalServerError
	}

	return fmt.Sprintf("User %d banned %d", userID, userToBanID), nil

}
