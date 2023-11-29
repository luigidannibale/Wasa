package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	//This takes the userID from parameters, if fails to convert error 400
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("message: 'could not convert userID'")
		return
	}

	//Delete the like of the user from the photo
	stream, s, err := rt.db.GetStream(userID)

	//Checks for DB-side errrors(404,500)
	if err != nil {
		if e := err.Error(); e == "UserNotFound" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("message : " + s)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("message : " + s)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stream)
}
