package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

//LastMod 21/11

// DELETE method, takes an userID and a photoID, and removes the like of the user from the photo
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	photoID, err := strconv.Atoi(ps.ByName("photoID"))
	//Takes photoID and checks if an error occurred
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("message: 'could not convert photoID'")
		return
	}
	//Takes userID and checks if an error occurred
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("message: 'could not convert userID'")
		return
	}

	//Delete the like of the user from the photo
	s, err := rt.db.Unlike(userID, photoID)

	//Checks for DB-side errrors(404,500)
	if err != nil {
		if e := err.Error(); e == "UserNotFound" || e == "PhotoNotFound" {
			w.WriteHeader(http.StatusNotFound)
			if e == "UserNotFound" {
				json.NewEncoder(w).Encode(userID)
			} else if e == "PhotoNotFound" {
				json.NewEncoder(w).Encode(photoID)
			}
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode("message : " + s)
}
