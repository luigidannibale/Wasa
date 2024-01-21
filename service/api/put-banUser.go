package api

import (
	"encoding/json"
	"errors"
	"strconv"

	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/database"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
			http.Error(w, MsgAuthNotFound+e.Error(), http.StatusUnauthorized)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgServerErrorUserID+e.Error(), http.StatusInternalServerError)
		}
		return
	}
	userIDparam, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, MsgConvertionErrorUserID+err.Error(), http.StatusBadRequest)
		return
	}
	if userIDauth != userIDparam {
		http.Error(w, MsgAuthNoMatch, http.StatusForbidden)
		return
	}
	userID := userIDauth

	// Takes the id of the user to ban, and validates it
	userToBanID, err := strconv.Atoi(r.URL.Query().Get("userToBanID"))
	if err != nil {
		http.Error(w, MsgConvertionErrorUserToBanID, http.StatusBadRequest)
		return
	}
	e = rt.db.VerifyUserId(userToBanID)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, "UserToBanID "+MsgNotFound, http.StatusNotFound)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgServerError+" while identifying userToBanID", http.StatusInternalServerError)
		}
		return
	}

	// Checks if the user is trying to ban himself
	if userID == userToBanID {
		http.Error(w, "The banner and banned can't have the same id", http.StatusForbidden)
		return
	}

	// Creates the ban that must be put in the DB
	var ban utils.Ban
	ban.BannerID = userID
	ban.BannedID = userToBanID

	// Puts the ban in the db
	s, err := rt.db.CreateBan(ban)

	// Checks for DB errrors
	if err != nil {
		if errors.Is(err, database.ErrAlreadyDone) {
			w.WriteHeader(http.StatusOK)
		}
		if errors.Is(err, database.ErrInternalServerError) {
			http.Error(w, MsgServerError+" while creating the ban "+s, http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	// Operation successful, creates an OK, or a CREATED response
	e = json.NewEncoder(w).Encode(s)
	if e != nil {
		http.Error(w, "Ban created but an error occurred while encoding the message "+e.Error(), http.StatusInternalServerError)
	}
}
