package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/utils"
)

//LastMod 21/11

// PUT method, takes an userID and a photoID, and puts the like of the user from the photo
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	photoID, err := strconv.Atoi(ps.ByName("photoID"))
	//Takes photoID and checks if an error occurred
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("message: 'could not convert photoID'")
		return
	}
	//Takes userID and checks if an error occurred
	userID, err := strconv.Atoi(r.URL.Query().Get("userID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("message: 'could not convert userID'")
		return
	}

	var like utils.Like
	like.UserID = userID

	//Delete the like of the user from the photo
	s, err := rt.db.Like(photoID, like)

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
		} else if e == "AlreadyLiked" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	json.NewEncoder(w).Encode("message : " + s)
	return

}
