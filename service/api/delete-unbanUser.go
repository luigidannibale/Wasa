package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	userIDauth, e := strconv.Atoi(r.Header.Get("Authorization"))
	if e != nil {
		http.Error(w, "Couldn't identify userId for authentication "+e.Error(), http.StatusUnauthorized)
		return
	}
	e = rt.db.VerifyUserId(userIDauth)
	if e != nil {
		http.Error(w, "The userID provided for authentication can't be found", http.StatusUnauthorized)
		return
	}
	userIDparam, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, "Could not convert the userID "+err.Error(), http.StatusBadRequest)
		return
	}
	if userIDauth != userIDparam {
		http.Error(w, "Authentication userID and parameter userID don't match", http.StatusForbidden)
		return
	}
	userID := userIDauth

	userToUnbanID, err := strconv.Atoi(ps.ByName("BannedID"))
	if err != nil {
		http.Error(w, "Could not convert the userToUnbanID", http.StatusBadRequest)
		return
	}

	if userID == userToUnbanID {
		http.Error(w, "The banner and banned can't have the same id", http.StatusForbidden)
		return
	}
	s, err := rt.db.DeleteBan(userID, userToUnbanID)

	if err != nil {
		if e := err.Error(); e == "UserNotFound" || e == "BannedNotFound" {
			http.Error(w, s, http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	e = json.NewEncoder(w).Encode(s)
	if e != nil {
		http.Error(w, s, http.StatusInternalServerError)
	}
}
