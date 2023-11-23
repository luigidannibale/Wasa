package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	//This takes the userID from parameters
	userID, err := strconv.Atoi(ps.ByName("userID"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("message : could not convert the userID")
		return
	}

	//This should get the body
	var photo utils.Photo
	json.NewDecoder(r.Body).Decode(&photo)

	err = photo.Validate()
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(photo)
		return
	}

	photoID, s, err := rt.db.UploadPhoto(userID, photo)
	//Checks for DB-side errrors(404,406,500)
	if err != nil {
		if e := err.Error(); e == "UserNotFound" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode("message : " + s)
		return
	}

	json.NewEncoder(w).Encode(photoID)
}
