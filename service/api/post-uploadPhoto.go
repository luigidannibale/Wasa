package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

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
		http.Error(w, MsgAuthNotFound+er.Error(), http.StatusUnauthorized)
		return
	}
	e := rt.db.VerifyUserId(userID)
	if e != nil {
		if errors.Is(e, database.ErrNotFound) {
			http.Error(w, MsgAuthNotFound+e.Error(), http.StatusUnauthorized)
		}
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgServerErrorUserID+e.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Takes caption from params and image from body
	caption := r.URL.Query().Get("caption")

	err := r.ParseMultipartForm(20 << 20) // 20 MB limit
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	imageFile, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error while retrieving image from request body", http.StatusBadRequest)
		return
	}
	defer imageFile.Close()

	image, e := io.ReadAll(imageFile)
	if e != nil {
		http.Error(w, "Couldn't convert the image "+e.Error(), http.StatusBadRequest)
		return
	}

	var photo utils.Photo
	photo.UserId = userID
	photo.Caption = caption
	photo.Image = base64.StdEncoding.EncodeToString(image)
	photo.UploadTimestamp = time.Now()
	if e = photo.Validate(); e != nil {
		http.Error(w, MsgValidationErrorPhoto+e.Error(), http.StatusBadRequest)
		return
	}

	// Creates the photo and gets the Id
	photoId, s, e := rt.db.CreatePhoto(photo)
	if e != nil {
		if errors.Is(e, database.ErrInternalServerError) {
			http.Error(w, MsgServerError+" while creating the comment "+s, http.StatusInternalServerError)
		}
		return
	}

	// Operation successful, creates a CREATED response
	w.WriteHeader(http.StatusCreated)
	e = json.NewEncoder(w).Encode(fmt.Sprintf("Photo uploaded successfully with id %d", photoId))
	if e != nil {
		http.Error(w, fmt.Sprintf("Photo uploaded successfully with id %d, but an error occurred while encoding the message ", photoId)+e.Error(), http.StatusInternalServerError)
	}
}
