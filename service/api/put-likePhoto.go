package api

import (
	"encoding/json"
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
		http.Error(w, "Couldn't identify userID for authentication", http.StatusUnauthorized)
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

	//Takes the photoID and validates it
	photoID, er := strconv.Atoi(ps.ByName("photoID"))
	if er != nil {
		http.Error(w, "Couldn't identify userId for authentication", http.StatusUnauthorized)
		return
	}
	_, s, e := rt.db.GetPhoto(photoID)
	if e != nil {
		if e == database.NotFound {
			http.Error(w, "Couldn't find the photo", http.StatusNotFound)
			return
		}
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	//Creates the like that has to be put in the db
	var like utils.Like
	like.UserID = userID
	like.PhotoID = photoID

	//Puts the like in the db
	s, e = rt.db.CreateLike(like)

	//Checks for DB errors
	if e != nil {
		switch e {
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
		http.Error(w, "Like created, but an error occurred while encoding the message", http.StatusInternalServerError)
	}
}