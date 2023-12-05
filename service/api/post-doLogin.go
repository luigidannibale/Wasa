package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/database"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	//Takes the username from the body
	var username map[string]string
	var e error
	e = json.NewDecoder(r.Body).Decode(&username)
	//Checks for errors
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

	//Creates a user with by username
	id, s, err := rt.db.CreateUserByUsername(u.Username)

	//Checks for DB errors
	if err != nil {
		switch err {
		case database.InternalServerError:
			http.Error(w, "An error occurred on the server "+s, http.StatusInternalServerError)
			return
		}
	}

	//Operation successful, creates an OK, or a CREATED response
	var message string
	if s == "Created" {
		w.WriteHeader(http.StatusCreated)
		message = "User created with id : "
	} else {
		w.WriteHeader(http.StatusOK)
		message = "User already existed, logged with id : "
	}
	e = json.NewEncoder(w).Encode(message + strconv.Itoa(id))
	if e != nil {
		http.Error(w, "User logged but couldn't encode the message, the id is "+strconv.Itoa(id), http.StatusInternalServerError)
	}
}
