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
		http.Error(w, "Couldn't identify userID for authentication "+er.Error(), http.StatusUnauthorized)
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

	// Takes the photoID and validates it
	photoID, er := strconv.Atoi(ps.ByName("photoID"))
	if er != nil {
		http.Error(w, "Couldn't identify userId for authentication "+er.Error(), http.StatusUnauthorized)
		return
	}
	_, s, e := rt.db.GetPhoto(photoID)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, s, http.StatusNotFound)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, "An error occurred on ther server while getting the photo "+s, http.StatusInternalServerError)
		}
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
			http.Error(w, "An error has occurred on the server while creating the like "+s, http.StatusInternalServerError)
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
