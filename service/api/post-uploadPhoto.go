package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	userID, er := strconv.Atoi(r.Header.Get("Authorization"))

	if er != nil {
		http.Error(w, "Couldn't identify userId for authentication", http.StatusUnauthorized)
		return
	}

	caption := r.URL.Query().Get("caption")
	var image string
	var e error
	e = json.NewDecoder(r.Body).Decode(&image)
	if e != nil {
		http.Error(w, "Couldn't convert the image", http.StatusBadRequest)
	}

	var photo utils.Photo
	photo.Caption = caption
	photo.Image = image
	photo.UploadTimestamp = utils.Now()
	if e = photo.Validate(); e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
	}

	photoId, s, e := rt.db.CreatePhoto(userID, photo)

	if e != nil {
		http.Error(w, s, http.StatusNotFound)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fmt.Sprintf("Photo uploaded successfully with id %d", photoId))
}
