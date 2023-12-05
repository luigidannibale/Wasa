package api

import (
	"encoding/json"
	"strconv"

	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/database"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		switch e {
		case database.NotFound:
			http.Error(w, "The userID provided for authentication can't be found", http.StatusUnauthorized)
			return
		case database.InternalServerError:
			http.Error(w, "An error occurred on ther server while identifying userID", http.StatusInternalServerError)
			return
		}
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

	//Takes the id of the user to follow, and validates it
	userToFollowID, err := strconv.Atoi(r.URL.Query().Get("userToFollowID"))
	if err != nil {
		http.Error(w, "Could not convert the userToFollowID", http.StatusBadRequest)
		return
	}
	e = rt.db.VerifyUserId(userToFollowID)
	if e != nil {
		http.Error(w, "The userToFollowID can't be found", http.StatusNotFound)
		return
	}

	//Checks if user is trying to follow himself
	if userID == userToFollowID {
		http.Error(w, "The follower and followed can't have the same id", http.StatusForbidden)
		return
	}

	//Creates the follow that must be put in the DB
	var follow utils.Follow
	follow.FollowerID = userID
	follow.FollowedID = userToFollowID

	//Puts the follow in the DB
	s, err := rt.db.CreateFollow(follow)

	//Checks for DB errrors
	if err != nil {
		switch err {
		case database.AlreadyDone:
			w.WriteHeader(http.StatusOK)
			break
		case database.InternalServerError:
			http.Error(w, "An error has occurred on the server "+s, http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	//Operation successful, creates an OK, or a CREATED response
	e = json.NewEncoder(w).Encode(s)
	if e != nil {
		http.Error(w, "Follow created but an error occurred while encoding the message ", http.StatusInternalServerError)
	}
}
