package database

import (
	"database/sql"
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
	"github.com/rickb777/date"
)

/*
Errors that can be returned: (InternalServerError)
*/
func (db *appdbimpl) CreateUser(user utils.User) (int, string, error) {
	var userID int
	userID, s, e := db.CreateUserByUsername(user.Username)
	if e != nil {
		return userID, s, e
	}
	user.Id = userID
	_, s, e = db.UpdateUser(user)
	return userID, s, e
}

/*
Errors that can be returned: (InternalServerError)
*/
func (db *appdbimpl) CreateUserByUsername(username string) (int, string, error) {
	var userID int
	u, s, e := db.GetUserByUsername(username)
	if e != nil {
		switch e {
		case ErrNotFound:
			err := db.c.QueryRow(`INSERT INTO Users(Username) VALUES (?) RETURNING Id`, username).Scan(&userID)
			if err == nil {
				return userID, "Created", nil
			}
			return userID, err.Error(), ErrInternalServerError
		case ErrInternalServerError:
			return userID, s, ErrInternalServerError
		}
	}
	return u.Id, "Logged", nil
}

/*
Errors that can be returned: (NotFound, InternalServerError, UsernameTaken)
*/
func (db *appdbimpl) UpdateUser(userNull utils.User) (utils.User, string, error) {

	user, _, e1 := db.GetUserByUsername(userNull.Username)
	if e1 == nil && userNull.Id != user.Id {
		return user, "This username is already used by someone else", ErrUsernameTaken
	}
	if len(userNull.Name) != 0 {
		user.Name = userNull.Name
	}
	if len(userNull.Surname) != 0 {
		user.Surname = userNull.Surname
	}
	if len(userNull.DateOfBirth.String()) != 0 {
		user.DateOfBirth = userNull.DateOfBirth
	}
	res, err := db.c.Exec(`	UPDATE Users
							SET Username = ? ,Name = ? ,Surname = ?,DateOfBirth = ?
							Where Id = ?`, user.Username, user.Name, user.Surname, user.DateOfBirth.String(), user.Id)

	if err != nil {
		return user, err.Error(), ErrInternalServerError
	}
	x, e := res.RowsAffected()
	if x == 0 {
		return user, e.Error(), ErrNotFound
	}
	if e != nil {
		return user, e.Error(), ErrInternalServerError
	}
	return user, "User updated successfully", nil

}

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetUser(userID int) (utils.User, string, error) {
	var user utils.User
	var name, surname, dateOfBirth sql.NullString
	e := db.c.QueryRow(`SELECT Id,Username,Name,Surname,DateOfBirth
						FROM Users
						WHERE Id = ?`, userID).Scan(&user.Id, &user.Username, &name, &surname, &dateOfBirth)

	if name.Valid {
		user.Name = name.String
	}
	if surname.Valid {
		user.Surname = surname.String
	}
	if dateOfBirth.Valid {
		d, er := date.AutoParse(dateOfBirth.String)
		if er != nil {
			return user, er.Error(), ErrInternalServerError
		}
		user.DateOfBirth = d
	}
	if e == nil {
		return user, "User found successfully", nil
	}
	if errors.Is(e, sql.ErrNoRows) {
		return user, "Couldn't find the user", ErrNotFound
	}
	return user, e.Error(), ErrInternalServerError
}

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetUserByUsername(username string) (utils.User, string, error) {
	var user utils.User
	var name, surname, dateOfBirth sql.NullString
	e := db.c.QueryRow(`SELECT Id,Username,Name,Surname,DateOfBirth
						FROM Users
						WHERE Username = ?`, username).Scan(&user.Id, &user.Username, &name, &surname, &dateOfBirth)
	if name.Valid {
		user.Name = name.String
	}
	if surname.Valid {
		user.Surname = surname.String
	}
	if dateOfBirth.Valid {
		d, er := date.AutoParse(dateOfBirth.String)
		if er != nil {
			return user, er.Error(), ErrInternalServerError
		}
		user.DateOfBirth = d
	}
	if e == nil {
		return user, "User found successfully", nil
	}
	if errors.Is(e, sql.ErrNoRows) {
		return user, "Couldn't find the user", ErrNotFound
	}
	return user, e.Error(), ErrInternalServerError
}

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) VerifyUserId(userId int) error {
	_, _, e := db.GetUser(userId)
	return e
}
