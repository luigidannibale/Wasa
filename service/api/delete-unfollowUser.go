package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

//LastMod 21/11

// DELETE method, takes an userID and a followedID, and removes the followed from user following
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	//This takes the userID from parameters, if fails to convert error 400
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("message: 'could not convert userID'")
		return
	}
	//This takes the followedID from parameters, if fails to convert error 400
	followedID, err := strconv.Atoi(ps.ByName("followedID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("message: 'could not convert followedID'")
		return
	}

	//Delete from following of userID the userToUnfollow, if one of them
	//isn't found error 404, if something else goes wrong errror 500
	s, err := rt.db.Unfollow(userID, followedID)

	//Checks for DB-side errrors(404,500)
	if err != nil {
		if e := err.Error(); e == "UserNotFound" || e == "FollowedNotFound" {
			w.WriteHeader(http.StatusNotFound)
			if e == "UserNotFound" {
				json.NewEncoder(w).Encode(userID)
			} else if e == "FollowedNotFound" {
				json.NewEncoder(w).Encode(followedID)
			}
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode("message : " + s)
	return

}
