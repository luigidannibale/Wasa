package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	//Taking the username
	var username map[string]string
	var e error
	e = json.NewDecoder(r.Body).Decode(&username)

	//Handling BadRequest
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		e = json.NewEncoder(w).Encode("Couldn't decode the username " + e.Error())
		if e != nil {
			http.Error(w, "Couldn't decode the username ", http.StatusBadRequest)
		}
		return
	}

	//Database query
	id, s, err := rt.db.CreateUserByUsername(username["username"])

	//Handling errors
	if err != nil {
		if err.Error() == "InternalServerError" {
			w.WriteHeader(http.StatusInternalServerError)
		}
		e = json.NewEncoder(w).Encode(s)
		if e != nil {
			http.Error(w, "Couldn't encode message", http.StatusInternalServerError)
		}
		return
	}
	var message string
	if s == "Created" {
		w.WriteHeader(http.StatusCreated)
		message = "User created with id : "
	} else {
		w.WriteHeader(http.StatusOK)
		message = "User already existed, logged with id : "
	}
	e = json.NewEncoder(w).Encode(message + strconv.Itoa(id))
	if e != nil {
		http.Error(w, "User created but couldn't encode the id, ths id is "+strconv.Itoa(id), http.StatusInternalServerError)
	}
}
