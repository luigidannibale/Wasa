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

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	/*Authentication part :
	- takes the userID from authentication
	- validates it
	*/
	userID, e := strconv.Atoi(r.Header.Get("Authorization"))
	if e != nil {
		http.Error(w, "Couldn't identify userId for authentication "+e.Error(), http.StatusUnauthorized)
		return
	}
	e = rt.db.VerifyUserId(userID)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, "The userID provided for authentication can't be found", http.StatusUnauthorized)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, "An error occurred on ther server while identifying userID", http.StatusInternalServerError)
		}
		return
	}

	//Takes the id of the photo, and validates it
	photoID, err := strconv.Atoi(ps.ByName("photoID"))
	if err != nil {
		http.Error(w, "Could not convert the photoID", http.StatusBadRequest)
		return
	}
	photo, s, e := rt.db.GetPhoto(photoID)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, s, http.StatusNotFound)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, "An error occurred while validating the photo "+e.Error(), http.StatusInternalServerError)
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
	//Creates the like that must be deleted from DB
	var like utils.Like
	like.PhotoID = photoID
	like.UserID = userID

	//Deletes the like from DB
	s, err = rt.db.DeleteLike(like)

	//Checks for DB errors
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, s, http.StatusNotFound)
		}
		if errors.Is(err, database.ErrInternalServerError) {
			http.Error(w, "An error occurred on the server while deleting the like"+s, http.StatusInternalServerError)
		}
		return
	}
	//Operation successful, creates an OK response
	w.WriteHeader(http.StatusOK)
	e = json.NewEncoder(w).Encode(s)
	if e != nil {
		http.Error(w, "Like deleted successfully but an error occurred while encoding the message", http.StatusInternalServerError)
	}
}
