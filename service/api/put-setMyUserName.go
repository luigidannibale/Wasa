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

	userID, e := strconv.Atoi(ps.ByName("userID"))
	username := r.URL.Query().Get("username")

	//Handling BadRequest
	if e != nil || username == "" {
		w.WriteHeader(http.StatusBadRequest)
		var errMessage string
		if e != nil {
			errMessage = "Couldn't decode userID " + e.Error()
		} else {
			errMessage = "Error taking the username"
		}
		e = json.NewEncoder(w).Encode(errMessage)
		if e != nil {
			http.Error(w, errMessage, http.StatusBadRequest)
		}
		return
	}
	var userToUpdate utils.User

	userToUpdate.Id = userID
	userToUpdate.Username = username

	user, s, err := rt.db.UpdateUser(userToUpdate)

	//Handling errors
	if err != nil {
		if err.Error() == "UsernameTaken" {
			w.WriteHeader(http.StatusConflict)
		}
		if err.Error() == "AlreadySo" {
			w.WriteHeader(http.StatusOK)
		}
		if err.Error() == "InternalServerError" {
			w.WriteHeader(http.StatusInternalServerError)
		} else if err.Error() == "NotFound" {
			w.WriteHeader(http.StatusNotFound)
		}
		e = json.NewEncoder(w).Encode(s)
		if e != nil {
			http.Error(w, "Couldn't encode error message : "+s, http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	e = json.NewEncoder(w).Encode("Username set successfully, new username : " + user.Username)
	if e != nil {
		http.Error(w, "Successful operation, the username has been updated, but an error occurred encoding the message", http.StatusInternalServerError)
	}
}
