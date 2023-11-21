package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"strconv"
	"io"
)


func (rt *_router) folloUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	
	//body, _ := io.readAll(r.Body)	
	//This should get the body
	var user User
	json.Unmarshal(r.Body, &user)
	
	//This takes the userID from parameters
	userID, err := strconv.Atoi(ps.ByName("userID"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if Users[userID] == nil {
		w.WriteHeader(http.StatusNotFound)
		return userID
	}
	if Users[user.id] == nil {
		w.WriteHeader(http.StatusNotFound)
		return body
	}
	// Insert the user into the followed 		
	// Users[userID].followed.append(user)


}