package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	photoID := ps.ByName("photoID")
	photo := getPhotofrombody()

	if err := validateDataByID(constants.photoID, photoID); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//This should be taken from db
	user_photos := []Photos{}
	photoID := len(user_photos) + 1
	user_photos.append(photo)
	pushToDB()
	json.NewEncoder(w).Encode(photoID)
}
