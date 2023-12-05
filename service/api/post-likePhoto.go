package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	userID, er := strconv.Atoi(r.Header.Get("Authorization"))
	if er != nil {
		http.Error(w, "Couldn't identify userId for authentication", http.StatusUnauthorized)
		return
	}
	_, s, e := rt.db.GetUser(userID)
	if e != nil {
		if e.Error() == "NotFound" {
			http.Error(w, "Couldn't find the user", http.StatusNotFound)
			return
		}
		http.Error(w, "An error occurred on the server "+e.Error(), http.StatusInternalServerError)
		return
	}
	photoID, er := strconv.Atoi(ps.ByName("photoID"))
	if er != nil {
		http.Error(w, "Couldn't identify userId for authentication", http.StatusUnauthorized)
		return
	}
	_, s, e = rt.db.GetPhoto(photoID)
	if e != nil {
		if e.Error() == "NotFound" {
			http.Error(w, "Couldn't find the photo", http.StatusNotFound)
			return
		}
		http.Error(w, "An error occurred on the server "+e.Error(), http.StatusInternalServerError)
		return
	}

	var like utils.Like
	like.UserID = userID
	like.PhotoID = photoID
	s, e = rt.db.CreateLike(like)

	if e != nil {
		if e.Error() == "NotFound" {
			http.Error(w, s, http.StatusNotFound)
			return
		}
		http.Error(w, s, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	e = json.NewEncoder(w).Encode(fmt.Sprintf("Photo uploaded successfully with id %d", photoID))
	if e != nil {
		http.Error(w, "An error occurred on the server "+e.Error(), http.StatusInternalServerError)
	}
}
