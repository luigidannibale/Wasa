package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/database"
)

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	/*Authentication part :
	- takes the userID from authorization
	- validates it
	*/
	userIDauth, e := strconv.Atoi(r.Header.Get("Authorization"))
	if e != nil {
		http.Error(w, MsgAuthNotFound+e.Error(), http.StatusUnauthorized)
		return
	}
	e = rt.db.VerifyUserId(userIDauth)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, MsgAuthNotFound+e.Error(), http.StatusUnauthorized)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgServerErrorUserID+e.Error(), http.StatusInternalServerError)
		}
		return
	}
	userID := userIDauth

	//  Takes the id of the photo to delete, and validates it
	photoID, err := strconv.Atoi(ps.ByName(ParamPhotoID))
	if err != nil {
		http.Error(w, MsgConvertionErrorPhotoID, http.StatusBadRequest)
		return
	}
	photo, s, e := rt.db.GetPhoto(photoID)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, s, http.StatusNotFound)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgValidationErrorPhoto+e.Error(), http.StatusInternalServerError)
		}
		return
	}

	//  Checks if the user is trying to delete a photo he didn't post
	if userID != photo.UserId {
		http.Error(w, "The user has not post this photo, nor can delete it", http.StatusForbidden)
		return
	}

	//  Deletes the photo from DB
	s, err = rt.db.DeletePhoto(photo)

	//  Checks for DB errors
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, s, http.StatusNotFound)
		}
		if errors.Is(err, database.ErrInternalServerError) {
			http.Error(w, MsgServerError+" while deleting the photo "+s, http.StatusInternalServerError)
		}
		return
	}

	//  Operation successful, creates an OK response
	w.WriteHeader(http.StatusOK)
	e = json.NewEncoder(w).Encode(s)
	if e != nil {
		http.Error(w, "Photo deleted successfully but an error occurred while encoding the message", http.StatusInternalServerError)
	}
}
