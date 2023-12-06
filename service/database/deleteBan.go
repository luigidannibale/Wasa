package database

import (
	"fmt"

	"github.com/luigidannibale/Wasa/service/utils"
)

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
