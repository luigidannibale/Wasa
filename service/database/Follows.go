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
func (db *appdbimpl) GetFollowedList(userID int) ([]utils.User, string, error) {
	var followed []utils.User
	rows, e := db.c.Query(`SELECT Id, Username, Name, Surname, DateOfBirth
						FROM Follows
						JOIN Users ON FollowedID = Id
						WHERE FollowerID = ?`, userID)
	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return followed, "Couldn't find any followed", ErrNotFound
		}
		return followed, e.Error(), ErrInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		var u utils.User
		err := rows.Scan(&u.Id, &u.Username, &u.Name, &u.Surname, &u.DateOfBirth)
		if err != nil {
			return followed, e.Error(), ErrInternalServerError
		}
		followed = append(followed, u)
	}

	return followed, "List of followed found successfully", nil

}

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetFollowersList(userID int) ([]utils.User, string, error) {
	var followers []utils.User
	rows, e := db.c.Query(`SELECT Id, Username, Name, Surname, DateOfBirth
						FROM Follows
						JOIN Users ON FollowerID = Id
						WHERE FollowedID = ?`, userID)
	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return followers, "Couldn't find any follower", ErrNotFound
		}
		return followers, e.Error(), ErrInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		var u utils.User
		err := rows.Scan(&u.Id, &u.Username, &u.Name, &u.Surname, &u.DateOfBirth)
		if err != nil {
			return followers, e.Error(), ErrInternalServerError
		}
		followers = append(followers, u)
	}

	return followers, "List of followers found successfully", nil

}

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetFollow(follow utils.Follow) (string, error) {

	e := db.c.QueryRow(`SELECT FollowerID, FollowedID
						FROM Follows
						WHERE FollowerID = ? AND FollowedID = ?`, follow.FollowerID, follow.FollowedID)
	if e.Err() == nil {
		return "Follow found successfully", nil
	}
	if errors.Is(e.Err(), sql.ErrNoRows) {
		return "Couldn't find the follow", ErrNotFound
	}
	return e.Err().Error(), ErrInternalServerError
}

/*
Errors that can be returned: (InternalServerError, AlreadyDone)
*/
func (db *appdbimpl) CreateFollow(follow utils.Follow) (string, error) {
	userID, userToFollowID := follow.FollowerID, follow.FollowedID
	res, err := db.c.Exec(`INSERT OR IGNORE INTO Follows(FollowerID,FollowedID) VALUES (?,?)`, userID, userToFollowID)

	if err != nil {
		return err.Error(), ErrInternalServerError
	}
	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return "Already following", ErrAlreadyDone
	}
	return fmt.Sprintf("User %d started following %d", userID, userToFollowID), nil
}

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) DeleteFollow(follow utils.Follow) (string, error) {
	userID, userToUnfollowID := follow.FollowerID, follow.FollowedID

	res, err := db.c.Exec(`DELETE FROM Follows WHERE FollowerID = ? AND FollowedID = ?`, userID, userToUnfollowID)
	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return fmt.Sprintf("Couldn't find the user to unfollow %d in the following of user %d", userID, userToUnfollowID), ErrNotFound
	}
	if err != nil {
		return err.Error(), ErrInternalServerError
	}

	return fmt.Sprintf("User %d is no longer following %d", userID, userToUnfollowID), nil
}
