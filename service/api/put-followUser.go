package api

import (
	"encoding/json"
	"errors"
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
	userIDparam, err := strconv.Atoi(ps.ByName(ParamUserID))
	if err != nil {
		http.Error(w, MsgConvertionErrorUserID+err.Error(), http.StatusBadRequest)
		return
	}
	if userIDauth != userIDparam {
		http.Error(w, MsgAuthNoMatch, http.StatusForbidden)
		return
	}
	userID := userIDauth

	// Takes the id of the user to follow, and validates it
	userToFollowID, err := strconv.Atoi(r.URL.Query().Get("userToFollowID"))
	if err != nil {
		http.Error(w, MsgConvertionErrorUserToFollowID+err.Error(), http.StatusBadRequest)
		return
	}
	e = rt.db.VerifyUserId(userToFollowID)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, "UserToFollowID "+MsgNotFound+e.Error(), http.StatusNotFound)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgServerError+" while identifying userToFollowID"+e.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Checks if user is trying to follow himself
	if userID == userToFollowID {
		http.Error(w, "The follower and followed can't have the same id", http.StatusForbidden)
		return
	}

	// Check if the user searched banned the one who is trying to search him
	e = rt.db.CheckBan(userToFollowID, userID)
	if e == nil {
		http.Error(w, "UserToFollowID "+MsgNotFound, http.StatusNotFound)
		return
	}
	if errors.Is(e, database.ErrInternalServerError) {
		http.Error(w, MsgServerError, http.StatusInternalServerError)
		return
	}

	// Creates the follow that must be put in the DB
	var follow utils.Follow
	follow.FollowerID = userID
	follow.FollowedID = userToFollowID

	// Puts the follow in the DB
	s, err := rt.db.CreateFollow(follow)

	// Checks for DB errrors
	if err != nil {
		if errors.Is(err, database.ErrAlreadyDone) {
			w.WriteHeader(http.StatusOK)
		}
		if errors.Is(err, database.ErrInternalServerError) {
			http.Error(w, MsgServerError+" while creating the follow "+s, http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	// Operation successful, creates an OK, or a CREATED response
	e = json.NewEncoder(w).Encode(s)
	if e != nil {
		http.Error(w, "Follow created but an error occurred while encoding the message "+e.Error(), http.StatusInternalServerError)
	}
}
