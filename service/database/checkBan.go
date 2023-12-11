package database

import (
	"github.com/luigidannibale/Wasa/service/utils"
)

func (db *appdbimpl) CheckBan(banner int, banned int) error {
	var b utils.Ban
	b.BannerID = banner
	b.BannedID = banned
	_, er := db.GetBan(b)
	return er
}
