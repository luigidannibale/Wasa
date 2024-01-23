package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/database"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	/*Authentication part :
	- takes userID from auth
	- validates it
	*/
	userID, er := strconv.Atoi(r.Header.Get("Authorization"))
	if er != nil {
		http.Error(w, MsgAuthNotFound+er.Error(), http.StatusUnauthorized)
		return
	}
	e := rt.db.VerifyUserId(userID)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, MsgAuthNotFound+e.Error(), http.StatusUnauthorized)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgServerErrorUserID+e.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Takes the photoID and validates it
	photoID, er := strconv.Atoi(ps.ByName(ParamPhotoID))
	if er != nil {
		http.Error(w, MsgConvertionErrorPhotoID+er.Error(), http.StatusUnauthorized)
		return
	}
	photo, s, e := rt.db.GetPhoto(photoID)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, s, http.StatusNotFound)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgValidationErrorPhoto+s, http.StatusInternalServerError)
		}
		return

	}

	// Check if the user that posted the searched photo banned the one who is trying to search it
	e = rt.db.CheckBan(photo.UserId, userID)
	if e == nil {
		http.Error(w, MsgNotFoundPhoto, http.StatusNotFound)
		return
	}
	if errors.Is(e, database.ErrInternalServerError) {
		http.Error(w, MsgServerError, http.StatusInternalServerError)
		return
	}

	// Creates the like that has to be put in the db
	var like utils.Like
	like.UserID = userID
	like.PhotoID = photoID

	// Puts the like in the db
	s, e = rt.db.CreateLike(like)

	// Checks for DB errors
	if e != nil {
		if errors.Is(e, database.ErrAlreadyDone) {
			w.WriteHeader(http.StatusOK)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgServerError+" while creating the like "+s, http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	// Operation successful, creates an OK, or a CREATED response
	e = json.NewEncoder(w).Encode(s)
	if e != nil {
		http.Error(w, "Like created, but an error occurred while encoding the message "+e.Error(), http.StatusInternalServerError)
	}
}
