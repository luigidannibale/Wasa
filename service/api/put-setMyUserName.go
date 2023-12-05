package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/database"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	u, s, e := rt.db.GetUser(userIDauth)
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

	//Takes the username and validates it
	username := r.URL.Query().Get("username")
	var userToUpdate utils.User

	userToUpdate.Id = userID
	userToUpdate.Username = username

	if e := userToUpdate.ValidateUsername(); e != nil {
		http.Error(w, "Username not valid "+e.Error(), http.StatusBadRequest)
		return
	}

	if u.Username == username {
		//Operation successful, creates a OK response
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode("Username was already set so")
		if e != nil {
			http.Error(w, "Successful operation, the username was already set so, but an error occurred encoding the message", http.StatusInternalServerError)
		}
	}

	//Update the user (with only the username)
	user, s, err := rt.db.UpdateUser(userToUpdate)

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

	//Operation successful, creates a CREATED response
	w.WriteHeader(http.StatusCreated)
	e = json.NewEncoder(w).Encode("Username set successfully, new username : " + user.Username)
	if e != nil {
		http.Error(w, "Successful operation, the username has been updated, but an error occurred encoding the message", http.StatusInternalServerError)
	}
}
