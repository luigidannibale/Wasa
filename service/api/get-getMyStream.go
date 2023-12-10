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

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	/*Authentication part :
	- takes the 2 userID (from auth and from params,)
	- validates 1 of them
	- checks if them are equal
	*/
	userIDauth, e := strconv.Atoi(r.Header.Get("Authorization"))
	if e != nil {
		http.Error(w, "Couldn't identify userId for authentication "+e.Error(), http.StatusUnauthorized)
		return
	}
	e = rt.db.VerifyUserId(userIDauth)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, "The userID provided for authentication can't be found", http.StatusUnauthorized)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, "An error occurred on ther server while identifying userID", http.StatusInternalServerError)
		}
		return
	}
	userIDparam, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, "Could not convert the userID "+err.Error(), http.StatusBadRequest)
		return
	}
	if userIDauth != userIDparam {
		http.Error(w, "Authentication userID and parameter userID don't match", http.StatusForbidden)
		return
	}
	userID := userIDauth

	users, s, err := rt.db.GetFollowedList(userID)
	// Checks for DB errors
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, s, http.StatusNotFound)
		}
		if errors.Is(err, database.ErrInternalServerError) {
			http.Error(w, "An error occurred on ther server while taking the user "+s, http.StatusInternalServerError)
		}
		return
	}

	var stream []utils.Photo
	for i := 0; i < len(users); i++ {
		photos, s, err := rt.db.GetPhotos(users[i].Id)
		if err != nil {
			http.Error(w, "An error occurred on ther server while taking the photos "+s, http.StatusInternalServerError)
		}
		for x := 0; x < len(photos); x++ {
			stream = append(stream, photos[x])
		}
	}

	stream, err = utils.SortStreamByLast(stream)
	if err != nil {

	}

	// Operation successful, creates an OK response
	w.WriteHeader(http.StatusOK)
	e = json.NewEncoder(w).Encode(stream)
	if e != nil {
		http.Error(w, "Operation successful but an error occured while returning the user "+e.Error(), http.StatusInternalServerError)
	}
}
