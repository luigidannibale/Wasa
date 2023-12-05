package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/database"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) updateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	/*Authentication part :
	- takes the 2 userID (from auth and from params,)
	- validates 1 of them
	- checks if them are equal
	*/
	userIDauth, e := strconv.Atoi(r.Header.Get("Authorization"))
	if e != nil {
		http.Error(w, "Couldn't identify userId for authentication", http.StatusUnauthorized)
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
	userIDparam, e := strconv.Atoi(ps.ByName("userID"))
	if e != nil {
		http.Error(w, "Couldn't decode userID "+e.Error(), http.StatusBadRequest)
		return
	}
	if userIDauth != userIDparam {
		http.Error(w, "Authentication userID and parameter userID don't match", http.StatusForbidden)
		return
	}
	userID := userIDauth
	//Takes the user from the body and validates it
	var user utils.User
	e = json.NewDecoder(r.Body).Decode(&user)
	if e != nil {
		http.Error(w, "Couldn't decode the user ", http.StatusBadRequest)
		return
	}

	//Checks if userId from auth is the same of the id of the given user
	if user.Id != userID {
		http.Error(w, "Authentication userID and the id of the given user don't match", http.StatusForbidden)
		return
	}

	//Updates the user
	user, s, err := rt.db.UpdateUser(user)

	//Checks for DB errors
	if err != nil {
		switch err {
		case database.UsernameTaken:
			http.Error(w, s, http.StatusConflict)
			return
		case database.InternalServerError:
			http.Error(w, "An error occurred on the server "+s, http.StatusInternalServerError)
			return
		}
	}

	//Operation successful, creates an OK response
	w.WriteHeader(http.StatusOK)
	e = json.NewEncoder(w).Encode(s)
	if e != nil {
		http.Error(w, "Successful operation, the user has been updated, but an error occurred encoding the message", http.StatusInternalServerError)
	}
}
