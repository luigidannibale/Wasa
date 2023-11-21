package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"errors"
)

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	
	userID := ps.ByName("userID")
	
	if err := validateDataByID(constants.userID, userID); err != nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//This should be taken from db
	followed_users := map[string]User {}

	var stream []Photo
	for f_u := followed_users[i] range followed_users{
		stream.append(getLastPhoto(f_u))		
	}
	sort(stream,by_last_uploaded)
	json.NewEncoder(w).Encode(stream)
}
