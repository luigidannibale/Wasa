package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) updateUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var user utils.User
	var e error
	e = json.NewDecoder(r.Body).Decode(&user)

	//Handling BadRequest
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		e = json.NewEncoder(w).Encode("message : Couldn't decode the User " + e.Error())
		if e != nil {
			http.Error(w, "Couldn't decode the user ", http.StatusBadRequest)
		}
		return
	}

	//Database query
	user, s, err := rt.db.UpdateUser(user)

	//Handling errors
	if err != nil {
		if err.Error() == "InternalServerError" {
			w.WriteHeader(http.StatusInternalServerError)
		} else if err.Error() == "NotFound" {
			w.WriteHeader(http.StatusNotFound)
		}
		e = json.NewEncoder(w).Encode("message : " + s)
		if e != nil {
			http.Error(w, "Couldn't encode message", http.StatusInternalServerError)
		}
		return
	}
	e = json.NewEncoder(w).Encode("message : " + s)
	if e != nil {
		http.Error(w, "Successful operation, the user has been updated, but an error occurred encoding the message", http.StatusInternalServerError)
	}
}
