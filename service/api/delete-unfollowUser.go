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

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	/*Authentication part :
	- takes the 2 userID (from auth and from params,)
	- validates 1 of them
	- checks if them are equal
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
	userIDparam, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, MsgConvertionErrorUserID+err.Error(), http.StatusBadRequest)
		return
	}
	if userIDauth != userIDparam {
		http.Error(w, MsgAuthNoMatch, http.StatusForbidden)
		return
	}
	userID := userIDauth

	// Takes the id of the user to unfollow, and validates it
	userToUnfollowID, err := strconv.Atoi(ps.ByName("followedID"))
	if err != nil {
		http.Error(w, MsgConvertionErrorFollowedID, http.StatusBadRequest)
		return
	}
	e = rt.db.VerifyUserId(userToUnfollowID)
	if e != nil {
		http.Error(w, "The user to unfollow can't be found", http.StatusNotFound)
		return
	}

	// Checks if user is trying to unfollow himself
	if userID == userToUnfollowID {
		http.Error(w, "The follower and followed can't have the same id", http.StatusForbidden)
		return
	}

	// Check if the user to unfollow banned the one who is trying to unfollow him
	e = rt.db.CheckBan(userToUnfollowID, userID)
	if e == nil {
		http.Error(w, "UserToUnfollow "+MsgNotFound, http.StatusNotFound)
		return
	}
	if errors.Is(e, database.ErrInternalServerError) {
		http.Error(w, MsgServerError, http.StatusInternalServerError)
		return
	}

	// Creates the follow that must be deleted from DB
	var follow utils.Follow
	follow.FollowerID = userID
	follow.FollowedID = userToUnfollowID

	// Deletes the follow from DB
	s, err := rt.db.DeleteFollow(follow)

	// Checks for DB errors
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, s, http.StatusNotFound)
		}
		if errors.Is(err, database.ErrInternalServerError) {
			http.Error(w, MsgServerError+" while deleting the follow"+s, http.StatusInternalServerError)
		}
		return
	}
	// Operation successful, creates an OK response
	w.WriteHeader(http.StatusOK)
	e = json.NewEncoder(w).Encode(s)
	if e != nil {
		http.Error(w, "Follow deleted successfully but an error occurred while encoding the message", http.StatusInternalServerError)
	}
}
