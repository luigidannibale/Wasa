package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

//LastMod 21/11

// DELETE method, takes a photoID and deletes the photo
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	//Takes the photoID and checks for error
	id, err := strconv.Atoi(ps.ByName("photoID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("message: 'could not convert photoID'")
		return
	}

	//Delete the photo, if it isn't found error 404,
	//if something else goes wrong errror 500
	s, err := rt.db.DeletePhoto(id)

	//Checks for errors that can occur DB-side (404,500)
	if err != nil {
		if err.Error() == "PhotoNotFound" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode("message : " + s)
}
