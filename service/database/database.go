/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/luigidannibale/Wasa/service/utils"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	// Operations on Users table ------------------------------

	/*
		Errors that can be returned: (InternalServerError)
	*/
	CreateUser(utils.User) (int, string, error)
	/*
		Errors that can be returned: (InternalServerError)
	*/
	CreateUserByUsername(string) (int, string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError, UsernameTaken)
	*/
	UpdateUser(utils.User) (utils.User, string, error)
	/*
		Errors that can be returned: (NotFound, InternalServerError)
	*/
	GetUser(int) (utils.User, string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	GetUserByUsername(string) (utils.User, string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	VerifyUserId(int) error

	// --------------------------------------------------------

	// Operations on Follows table ----------------------------

	/*
		Errors that can be returned: (NotFound, InternalServerError)
	*/
	GetFollowedList(int) ([]utils.User, string, error)
	/*
		Errors that can be returned: (NotFound, InternalServerError)
	*/
	GetFollowersList(int) ([]utils.User, string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	GetFollow(utils.Follow) (string, error)
	/*
		Errors that can be returned: (InternalServerError, AlreadyDone)
	*/
	CreateFollow(utils.Follow) (string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	DeleteFollow(utils.Follow) (string, error)

	// --------------------------------------------------------

	// Operations on Bans table -------------------------------

	/*
		Errors that can be returned: (NotFound, InternalServerError)
	*/
	GetBannedList(int) ([]utils.User, string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	GetBan(utils.Ban) (string, error)
	/*
		Errors that can be returned: (InternalServerError, AlreadyDone)
	*/
	CreateBan(utils.Ban) (string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	CheckBan(int, int) error
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	DeleteBan(utils.Ban) (string, error)

	// --------------------------------------------------------

	// Operations on Photos table -----------------------------

	/*
		Errors that can be returned: (InternalServerError)
	*/
	CreatePhoto(utils.Photo) (int, string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	GetPhoto(int) (utils.Photo, string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	GetStream(int) ([]utils.Photo, string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	GetMyPhotos(int) ([]utils.Photo, string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	DeletePhoto(utils.Photo) (string, error)

	// --------------------------------------------------------

	// Operations on Likes table ------------------------------
	/*
	   Errors that can be returned: (InternalServerError, AlreadyDone)
	*/
	CreateLike(utils.Like) (string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	GetLike(utils.Like) (string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	DeleteLike(utils.Like) (string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	GetLikersList(int) ([]string, string, error)

	// --------------------------------------------------------

	// Operations on Comment table ----------------------------
	/*
		Errors that can be returned: (InternalServerError)
	*/
	CreateComment(utils.Comment) (int, string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	GetComment(int) (utils.Comment, string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	DeleteComment(int) (string, error)
	/*
	   Errors that can be returned: (NotFound, InternalServerError)
	*/
	GetCommentsList(int) ([]utils.Comment, string, error)

	// --------------------------------------------------------
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	var tableName string
	// Creates Users table if not already existing #LastMod - 10/12
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' and name = 'Users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = CreateUsersTable(db)
		if err != nil {
			return nil, fmt.Errorf("error creating 'Users' table: %w", err)
		}
	}
	// Creates Photos table if not already existing #LastMod - 10/12
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' and name = 'Photos';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = CreatePhotosTable(db)
		if err != nil {
			return nil, fmt.Errorf("error creating 'Photos' table: %w", err)
		}
	}
	// Creates Comments table if not already existing #LastMod - 10/12
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' and name = 'Comments';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = CreateCommentsTable(db)
		if err != nil {
			return nil, fmt.Errorf("error creating 'Comments' table: %w", err)
		}
	}
	// Creates Likes table if not already existing #LastMod - 10/12
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' and name = 'Likes';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = CreateLikesTable(db)
		if err != nil {
			return nil, fmt.Errorf("error creating 'Likes' table: %w", err)
		}
	}
	// Creates Follows table if not already existing #LastMod - 10/12
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' and name = 'Follows';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = CreateFollowsTable(db)
		if err != nil {
			return nil, fmt.Errorf("error creating 'Follows' table: %w", err)
		}
	}
	// Creates Bans table if not already existing #LastMod - 10/12
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' and name = 'Bans';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = CreateBansTable(db)
		if err != nil {
			return nil, fmt.Errorf("error creating 'Bans' table: %w", err)
		}
	}
	// Populates the DB with some values
	err = PopulateDB(db)
	if err != nil {
		return nil, fmt.Errorf("error populating db: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

func CreateUsersTable(db *sql.DB) error {
	var sqlStmt string
	var err error
	sqlStmt = `CREATE TABLE Users
			(Id INTEGER NOT NULL PRIMARY KEY  UNIQUE, 
			Username TEXT NOT NULL UNIQUE,
			Name TEXT,
			Surname TEXT,
			DateOfBirth TEXT);`
	_, err = db.Exec(sqlStmt)
	return err
}
func CreatePhotosTable(db *sql.DB) error {
	var sqlStmt string
	var err error
	sqlStmt = `CREATE TABLE Photos 
			(Id INTEGER NOT NULL PRIMARY KEY  UNIQUE, 
			UserID INTEGER,
			Image BLOB NOT NULL,
			Caption TEXT,
			UploadTimestamp TEXT NOT NULL,
			FOREIGN KEY(UserID) REFERENCES Users(Id));`
	_, err = db.Exec(sqlStmt)
	return err
}
func CreateCommentsTable(db *sql.DB) error {
	var sqlStmt string
	var err error
	sqlStmt = `CREATE TABLE Comments 
			(Id INTEGER NOT NULL PRIMARY KEY  UNIQUE, 
			UserID INTEGER,
			PhotoID INTEGER,
			Content TEXT NOT NULL,
			FOREIGN KEY(PhotoID) REFERENCES Photos(Id) ON DELETE CASCADE,
			FOREIGN KEY(UserID) REFERENCES Users(Id));`
	_, err = db.Exec(sqlStmt)
	return err
}
func CreateLikesTable(db *sql.DB) error {
	var sqlStmt string
	var err error
	sqlStmt = `CREATE TABLE Likes 
			(UserID INTEGER,
			PhotoID INTEGER,
			PRIMARY KEY(UserID, PhotoID),
			FOREIGN KEY(PhotoID) REFERENCES Photos(Id) ON DELETE CASCADE,
			FOREIGN KEY(UserID) REFERENCES Users(Id));`
	_, err = db.Exec(sqlStmt)
	return err
}
func CreateFollowsTable(db *sql.DB) error {
	var sqlStmt string
	var err error
	sqlStmt = `CREATE TABLE Follows 
			(FollowerID INTEGER,
			FollowedID INTEGER,
			PRIMARY KEY(FollowerID, FollowedID),
			FOREIGN KEY(FollowerID) REFERENCES Users(Id),
			FOREIGN KEY(FollowedID) REFERENCES Users(Id));`
	_, err = db.Exec(sqlStmt)
	return err
}
func CreateBansTable(db *sql.DB) error {
	var sqlStmt string
	var err error
	sqlStmt = `CREATE TABLE Bans 
			(BannerID INTEGER,
			BannedID INTEGER,
			PRIMARY KEY(BannerID, BannedID),
			FOREIGN KEY(BannerID) REFERENCES Users(Id),
			FOREIGN KEY(BannedID) REFERENCES Users(Id));`
	_, err = db.Exec(sqlStmt)
	return err
}
func PopulateDB(db *sql.DB) error {
	var sqlStmt string
	var err error

	sqlStmt = `INSERT INTO Users (Username,Name,Surname,DateOfBirth)
				VALUES
						("Gigi","Luigi","Dannibale","2002-12-16"),
						("Paoletto","Paolo","Rossi","2002-12-23"),
						("Gianni","Gianfranco","Verdi","2002-11-19"),
						("Paolino","Paolo","Bianchi","2002-01-15"),
						("Getto","Giorgio","Verdi","2003-01-15"),
						("Matto","Matteo","Bianchi","2002-01-15"),
						("Fra","Francesco","Crema","2004-08-19");`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error populating users table: %w", err)
	}

	sqlStmt = `INSERT INTO Follows (FollowerID,FollowedID)
				VALUES
						(1,2),
						(1,4),
						(2,1),						
						(2,4),
						(3,5),
						(3,4);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error populating follows table: %w", err)
	}

	sqlStmt = `INSERT INTO Bans (BannerID,BannedID)
				VALUES
						(1,5),
						(1,6),
						(2,3),
						(3,1),
						(3,2),
						(5,6);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error populating bans table: %w", err)
	}

	sqlStmt = `INSERT INTO Photos (UserID,Image,Caption,UploadTimestamp)
				VALUES
						(1,"","First photo","11/11 03:08:38PM '22 +0100"),
						(2,"","Second photo","12/12 02:08:38PM '22 +0100"),
						(3,"","Third photo","11/11 04:08:38PM '23 +0100"),
						(4,"","Fourth photo","12/12 05:08:38PM '23 +0100"),
						(5,"","Fifth photo","13/12 07:08:38PM '23 +0100"),
						(6,"","Sixth photo","14/12 01:08:38PM '23 +0100");`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error populating photos table: %w", err)
	}

	sqlStmt = `INSERT INTO Comments (UserID,PhotoID,Content)
				VALUES
						(1,1,"this is my comment to my photo"),
						(1,2,"nice pic"),
						(2,2,"this is my comment to my photo"),
						(2,1,"wow"),
						(3,3,"this is my comment to my photo"),
						(4,4,"this is my comment to my photo"),
						(5,5,"this is my comment to my photo"),
						(6,6,"this is my comment to my photo");`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error populating comments table: %w", err)
	}

	sqlStmt = `INSERT INTO Likes (UserID,PhotoID)
				VALUES
						(1,2),
						(1,3),
						(2,1);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error populating likes table: %w", err)
	}

	return nil
}
