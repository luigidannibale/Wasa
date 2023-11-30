package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	//This should get the body
	var user utils.User
	var e error
	e = json.NewDecoder(r.Body).Decode(&user)
	//Handling BadRequest
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		e = json.NewEncoder(w).Encode("message : Couldn't decode the User")
		if e != nil {
			http.Error(w, "Couldn't decode the user", http.StatusBadRequest)
		}
		return
	}

	//Database query
	id, s, err := rt.db.CreateUser(user.Username)

	//Handling errors
	if err != nil {
		if err.Error() == "InternalServerError" {
			w.WriteHeader(http.StatusInternalServerError)
		}
		e = json.NewEncoder(w).Encode("message : " + s)
		if e != nil {
			http.Error(w, "Couldn't encode message", http.StatusInternalServerError)
		}
		return
	}
	if s == "Created" {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	e = json.NewEncoder(w).Encode(id)
	if e != nil {
		http.Error(w, "User created but couldn't ecnode the id, ths id is "+strconv.Itoa(id), http.StatusInternalServerError)
	}
}
