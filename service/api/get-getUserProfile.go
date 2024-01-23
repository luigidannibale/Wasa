package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/database"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	/*Authentication part :
	- takes userID from auth
	- validates it
	*/
	userID, er := strconv.Atoi(r.Header.Get("Authorization"))
	if er != nil {
		http.Error(w, MsgAuthNotFound+er.Error(), http.StatusUnauthorized)
		return
	}
	e := rt.db.VerifyUserId(userID)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, MsgAuthNotFound+e.Error(), http.StatusUnauthorized)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgServerErrorUserID+e.Error(), http.StatusInternalServerError)
		}
		return
	}
	// Takes the username from params and validates it
	username := r.URL.Query().Get("username")
	var u utils.User
	u.Username = username

	if e := u.ValidateUsername(); e != nil {
		http.Error(w, "Username not valid "+e.Error(), http.StatusBadRequest)
		return
	}

	// Gets the user by the username
	user, s, err := rt.db.GetUserByUsername(username)

	// Checks for DB errors
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, s, http.StatusNotFound)
		}
		if errors.Is(err, database.ErrInternalServerError) {
			http.Error(w, MsgServerError+" while taking the user"+s, http.StatusInternalServerError)
		}
		return
	}

	// Check if the user searched banned the one who is trying to search him
	e = rt.db.CheckBan(user.Id, userID)
	if e == nil {
		http.Error(w, MsgNotFoundUserID, http.StatusNotFound)
		return
	}
	if errors.Is(e, database.ErrInternalServerError) {
		http.Error(w, MsgServerError, http.StatusInternalServerError)
		return
	}

	// Operation successful, creates an OK response
	w.WriteHeader(http.StatusOK)
	e = json.NewEncoder(w).Encode(user)
	if e != nil {
		http.Error(w, "Operation successful but an error occured while returning the user "+e.Error(), http.StatusInternalServerError)
	}
}
