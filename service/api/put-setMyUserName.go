package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	userIDauth, e := strconv.Atoi(r.Header.Get("Authorization"))
	if e != nil {
		http.Error(w, "Couldn't identify userId for authentication", http.StatusUnauthorized)
		return
	}
	e = rt.db.VerifyUserId(userIDauth)
	if e != nil {
		http.Error(w, "The userID provided for authentication can't be found", http.StatusUnauthorized)
		return
	}
	userIDparam, e := strconv.Atoi(ps.ByName("userID"))
	if e != nil {
		http.Error(w, "Couldn't decode userID "+e.Error(), http.StatusBadRequest)
		return
	}
	if userIDauth != userIDparam {
		http.Error(w, "Authentication userID and parameter userID don't match", http.StatusForbidden)
		return
	}
	userID := userIDauth
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Error taking the username", http.StatusBadRequest)
		return
	}

	var userToUpdate utils.User

	userToUpdate.Id = userID
	userToUpdate.Username = username

	user, s, err := rt.db.UpdateUser(userToUpdate)

	if err != nil {
		if err.Error() == "UsernameTaken" {
			http.Error(w, s, http.StatusConflict)
		}
		if err.Error() == "AlreadySo" {
			w.WriteHeader(http.StatusOK)
		}
		if err.Error() == "InternalServerError" {
			http.Error(w, s, http.StatusInternalServerError)
		}
		e = json.NewEncoder(w).Encode(s)
		if e != nil {
			http.Error(w, s, http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	e = json.NewEncoder(w).Encode("Username set successfully, new username : " + user.Username)
	if e != nil {
		http.Error(w, "Successful operation, the username has been updated, but an error occurred encoding the message", http.StatusInternalServerError)
	}
}
