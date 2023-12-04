package api

import (
	"encoding/json"
	"strconv"

	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	userToBanID, err := strconv.Atoi(r.URL.Query().Get("userToBanID"))
	if err != nil {
		http.Error(w, "Could not convert the userToBanID", http.StatusBadRequest)
		return
	}
	e = rt.db.VerifyUserId(userToBanID)
	if e != nil {
		http.Error(w, "The userToBanID can't be found", http.StatusNotFound)
		return
	}

	if userID == userToBanID {
		http.Error(w, "The banner and banned can't have the same id", http.StatusForbidden)
		return
	}
	s, err := rt.db.CreateBan(userID, userToBanID)

	//Checks for DB-side errrors(404,500)
	if err != nil {
		if err.Error() == "AlreadyBanned" {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, s, http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	e = json.NewEncoder(w).Encode(s)
	if e != nil {
		http.Error(w, s, http.StatusInternalServerError)
	}
}
