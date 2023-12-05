package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/database"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		switch e {
		case database.NotFound:
			http.Error(w, "The userID provided for authentication can't be found", http.StatusUnauthorized)
			return
		case database.InternalServerError:
			http.Error(w, "An error occurred on ther server while identifying userID", http.StatusInternalServerError)
			return
		}
	}
	//Takes the username from params and validates it
	username := r.URL.Query().Get("username")
	var u utils.User
	u.Username = username

	if e := u.ValidateUsername(); e != nil {
		http.Error(w, "Username not valid "+e.Error(), http.StatusBadRequest)
		return
	}

	//Gets the user by the username
	user, s, err := rt.db.GetUserByUsername(username)

	//Checks for DB errors
	if err != nil {
		switch err {
		case database.NotFound:
			http.Error(w, s, http.StatusNotFound)
			return
		case database.InternalServerError:
			http.Error(w, "An error occurred on ther server "+s, http.StatusInternalServerError)
			return
		}
	}

	//Operation successful, creates an OK response
	w.WriteHeader(http.StatusOK)
	e = json.NewEncoder(w).Encode(user)
	if e != nil {
		http.Error(w, "Operation successful but an error occured while returning the user ", http.StatusInternalServerError)
	}
}
