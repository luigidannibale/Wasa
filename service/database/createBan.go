package database

import (
	"errors"
	"fmt"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (InternalServerError, AlreadyDone)
*/
func (db *appdbimpl) CreateBan(ban utils.Ban) (string, error) {
	userID, userToBanID := ban.BannerID, ban.BannedID
	res, err := db.c.Exec(`INSERT OR IGNORE INTO Bans(BannerID,BannedID)
									VALUES (?,?)`, userID, userToBanID)
	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return "Already banned", ErrAlreadyDone
	}
	if err != nil {
		return err.Error(), errors.New("InternalServerError")
	}
	return fmt.Sprintf("User %d banned %d", userID, userToBanID), nil

}
