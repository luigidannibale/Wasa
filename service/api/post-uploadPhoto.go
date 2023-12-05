package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/database"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	//Takes caption from params and image from body
	caption := r.URL.Query().Get("caption")
	var image string
	e = json.NewDecoder(r.Body).Decode(&image)
	if e != nil {
		http.Error(w, "Couldn't convert the image", http.StatusBadRequest)
	}

	var photo utils.Photo
	photo.UserId = userID
	photo.Caption = caption
	photo.Image = image
	photo.UploadTimestamp = utils.Now()
	if e = photo.Validate(); e != nil {
		http.Error(w, "Couldn't validate the photo "+e.Error(), http.StatusBadRequest)
	}

	//Creates the photo and gets the Id
	photoId, s, e := rt.db.CreatePhoto(photo)
	if e != nil {
		switch e {
		case database.InternalServerError:
			http.Error(w, "An error occurred on the server "+s, http.StatusInternalServerError)
			return
		}
	}

	//Operation successful, creates a CREATED response
	w.WriteHeader(http.StatusCreated)
	e = json.NewEncoder(w).Encode(fmt.Sprintf("Photo uploaded successfully with id %d", photoId))
	if e != nil {
		http.Error(w, fmt.Sprintf("Photo uploaded successfully with id %d, but an error occurred while encoding the message ", photoId)+e.Error(), http.StatusInternalServerError)
	}
}
