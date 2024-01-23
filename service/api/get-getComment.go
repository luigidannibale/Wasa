package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/database"
)

func (rt *_router) getComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// Takes the photoID from params and validates it
	photoID, e := strconv.Atoi(ps.ByName(ParamPhotoID))
	if e != nil {
		http.Error(w, MsgConvertionErrorPhotoID+e.Error(), http.StatusBadRequest)
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

	// Takes the commentID from params and validates it
	commentID, e := strconv.Atoi(ps.ByName("commentID"))
	if e != nil {
		http.Error(w, MsgConvertionErrorCommentID+e.Error(), http.StatusBadRequest)
		return
	}

	// Gets the comment
	comment, s, err := rt.db.GetComment(commentID)

	// Checks for DB errors
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, s, http.StatusNotFound)
		}
		if errors.Is(err, database.ErrInternalServerError) {
			http.Error(w, MsgServerError+" while getting the comment"+s, http.StatusInternalServerError)
		}
		return
	}

	// Checks if the comment is of another photo
	if comment.PhotoID != photoID {
		http.Error(w, "There is no comment to this photo with such id", http.StatusNotFound)
		return
	}

	// Operation successful, creates an OK response
	w.WriteHeader(http.StatusOK)
	e = json.NewEncoder(w).Encode(comment)
	if e != nil {
		http.Error(w, "Operation successful but an error occured while returning the comment "+e.Error(), http.StatusInternalServerError)
	}
}
