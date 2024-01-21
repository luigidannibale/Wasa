package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/database"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Takes the username from the body
	var username map[string]string
	var e error
	e = json.NewDecoder(r.Body).Decode(&username)
	// Checks for errors
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		e = json.NewEncoder(w).Encode("Couldn't decode the username " + e.Error())
		if e != nil {
			http.Error(w, "Couldn't decode the username ", http.StatusBadRequest)
		}
		return
	}
	var u utils.User
	u.Username = username["username"]

	if e := u.ValidateUsername(); e != nil {
		http.Error(w, "Username not valid "+e.Error(), http.StatusBadRequest)
		return
	}

	// Creates a user with by username
	id, s, err := rt.db.CreateUserByUsername(u.Username)

	// Checks for DB errors
	if err != nil {
		if errors.Is(err, database.ErrInternalServerError) {
			http.Error(w, MsgServerError+" while creating the user "+s, http.StatusInternalServerError)
		}
		return
	}

	// Operation successful, creates an OK, or a CREATED response
	var message string
	if s == "Created" {
		w.WriteHeader(http.StatusCreated)
		message = fmt.Sprintf("User created with id : %d", id)
	} else {
		w.WriteHeader(http.StatusOK)
		message = fmt.Sprintf("User already existed, logged with id : %d", id)
	}
	e = json.NewEncoder(w).Encode(message)
	if e != nil {
		http.Error(w, fmt.Sprintf("User logged with id %d, but an error occurred while encoding the message ", id)+e.Error(), http.StatusInternalServerError)
	}
}
