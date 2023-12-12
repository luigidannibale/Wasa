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
		http.Error(w, "Couldn't identify userId for authentication", http.StatusUnauthorized)
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

	// Takes the commentID from params and validates it
	commentID, e := strconv.Atoi(ps.ByName("commentID"))
	if e != nil {
		http.Error(w, "Error taking the commentID "+e.Error(), http.StatusBadRequest)
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
			http.Error(w, "An error occurred on ther server while getting the comment"+s, http.StatusInternalServerError)
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
