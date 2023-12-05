package database

func (db *appdbimpl) VerifyUserId(userId int) error {
	_, _, e := db.GetUser(userId)
	return e
}
