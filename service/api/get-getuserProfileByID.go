package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/database"
)

func (rt *_router) getUserProfileByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	fmt.Println("Ciao")
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
			http.Error(w, MsgAuthNotFound+e.Error(), http.StatusUnauthorized)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgServerErrorUserID+e.Error(), http.StatusInternalServerError)
		}
		return
	}
	userIDparam, err := strconv.Atoi(ps.ByName(ParamUserID))
	if err != nil {
		http.Error(w, MsgConvertionErrorUserID+err.Error(), http.StatusBadRequest)
		return
	}
	e = rt.db.VerifyUserId(userIDparam)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, MsgAuthNotFound+e.Error(), http.StatusNotFound)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgServerErrorUserID+e.Error(), http.StatusInternalServerError)
		}
		return
	}
	userID := userIDparam

	// Gets the user by the username
	user, s, err := rt.db.GetUser(userID)

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
