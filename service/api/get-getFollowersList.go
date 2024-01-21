package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/database"
)

func (rt *_router) getFollowersList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	/*Authentication part :
	- takes the 2 userID (from auth and from params,)
	- validates 1 of them
	- checks if them are equal
	*/
	userIDauth, e := strconv.Atoi(r.Header.Get("Authorization"))
	if e != nil {
		http.Error(w, MsgAuthNotFound+e.Error(), http.StatusUnauthorized)
		return
	}
	e = rt.db.VerifyUserId(userIDauth)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, MsgConvertionErrorUserID+e.Error(), http.StatusNotFound)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgConvertionErrorUserID+e.Error(), http.StatusInternalServerError)
		}
		return
	}
	userIDparam, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, MsgConvertionErrorUserID+err.Error(), http.StatusBadRequest)
		return
	}
	e = rt.db.VerifyUserId(userIDparam)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, "UserID "+MsgNotFound+e.Error(), http.StatusNotFound)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgServerErrorUserID+e.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if the user searched banned the one who is trying to search him
	if userIDauth != userIDparam {
		e = rt.db.CheckBan(userIDparam, userIDauth)
		if e == nil {
			http.Error(w, "UserID "+MsgNotFound, http.StatusNotFound)
			return
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgServerErrorUserID, http.StatusInternalServerError)
			return
		}
	}

	// Gets the list of followed
	followersList, s, err := rt.db.GetFollowersList(userIDparam)

	// Checks for DB errors
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, s, http.StatusNotFound)
		}
		if errors.Is(err, database.ErrInternalServerError) {
			http.Error(w, MsgServerError+" while getting the list of followers"+s, http.StatusInternalServerError)
		}
		return
	}

	// Operation successful, creates an OK response
	w.WriteHeader(http.StatusOK)
	e = json.NewEncoder(w).Encode(followersList)
	if e != nil {
		http.Error(w, "Operation successful but an error occured while returning the list of followed "+e.Error(), http.StatusInternalServerError)
	}
}
