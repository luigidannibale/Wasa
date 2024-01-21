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

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	//  Takes the id of the user to unban, and validates it
	userToUnbanID, err := strconv.Atoi(ps.ByName("bannedID"))
	if err != nil {
		http.Error(w, MsgConvertionErrorBannedID, http.StatusBadRequest)
		return
	}
	//  Checks if the user is trying to unban himself
	if userID == userToUnbanID {
		http.Error(w, "The banner and banned can't have the same id", http.StatusForbidden)
		return
	}
	e = rt.db.VerifyUserId(userToUnbanID)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, "The user to unban can't be found", http.StatusNotFound)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgServerError+" while looking for the user to unban ", http.StatusInternalServerError)
		}
		return
	}

	//  Creates the ban that must be deleted from DB
	var ban utils.Ban
	ban.BannerID = userID
	ban.BannedID = userToUnbanID

	//  Deletes the ban from DB
	s, err := rt.db.DeleteBan(ban)

	//  Checks for DB errors
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, s, http.StatusNotFound)
		}
		if errors.Is(err, database.ErrInternalServerError) {
			http.Error(w, MsgServerError+" while deleting the ban"+s, http.StatusInternalServerError)
		}
		return
	}

	//  Operation successful, creates an OK response
	w.WriteHeader(http.StatusOK)
	e = json.NewEncoder(w).Encode(s)
	if e != nil {
		http.Error(w, "Ban deleted successfully but an error occurred while encoding the message", http.StatusInternalServerError)
	}
}
