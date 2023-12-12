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
		http.Error(w, "Couldn't identify userId for authentication "+er.Error(), http.StatusUnauthorized)
		return
	}
	e := rt.db.VerifyUserId(userID)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, "The userID provided for authentication can't be found", http.StatusUnauthorized)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, "An error occurred on ther server while identifying userID", http.StatusInternalServerError)
		}
		return
	}
	// Takes the photoID from params and validates it
	photoID, e := strconv.Atoi(ps.ByName("photoID"))
	if e != nil {
		http.Error(w, "Error taking the photoID "+e.Error(), http.StatusBadRequest)
		return
	}
	photo, s, e := rt.db.GetPhoto(photoID)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, s, http.StatusNotFound)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, "An error occurred while validating the photo "+s, http.StatusInternalServerError)
		}
		return
	}

	// Check if the user that posted the searched photo banned the one who is trying to search it
	e = rt.db.CheckBan(photo.UserId, userID)
	if e == nil {
		http.Error(w, "Couldn't find the photo", http.StatusNotFound)
		return
	}
	if errors.Is(e, database.ErrInternalServerError) {
		http.Error(w, "An error occurred on ther server", http.StatusInternalServerError)
		return
	}

	// Takes content from params
	content := r.URL.Query().Get("content")

	var comment utils.Comment
	comment.PhotoID = photoID
	comment.UserID = userID
	comment.Content = content

	if e = comment.Validate(); e != nil {
		http.Error(w, "Couldn't validate the photo "+e.Error(), http.StatusBadRequest)
		return
	}

	// Creates the comment and gets the Id
	commentID, s, e := rt.db.CreateComment(comment)
	if e != nil {
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, "An error occurred on the server creating the comment "+s, http.StatusInternalServerError)
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
