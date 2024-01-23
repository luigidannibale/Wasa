package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/database"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// Takes content from params
	content := r.URL.Query().Get("content")

	var comment utils.Comment
	comment.PhotoID = photoID
	comment.UserID = userID
	comment.Content = content

	if e = comment.Validate(); e != nil {
		http.Error(w, "Couldn't validate the comment "+e.Error(), http.StatusBadRequest)
		return
	}

	// Creates the comment and gets the Id
	commentID, s, e := rt.db.CreateComment(comment)
	if e != nil {
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgServerError+" while creating the comment "+s, http.StatusInternalServerError)
		}
		return
	}

	// Operation successful, creates a CREATED response
	w.WriteHeader(http.StatusCreated)
	e = json.NewEncoder(w).Encode(fmt.Sprintf("Comment posted successfully with id %d", commentID))
	if e != nil {
		http.Error(w, fmt.Sprintf("Comment posted successfully with id %d, but an error occurred while encoding the message ", commentID)+e.Error(), http.StatusInternalServerError)
	}
}
