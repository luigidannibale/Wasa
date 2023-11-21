package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	//This takes the userID from parameters
	userID, err := strconv.Atoi(ps.ByName("userID"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("message : could not convert the userID")
		return
	}
	body, _ := io.ReadAll(r.Body)
	//This should get the body
	var userToFollow utils.User
	json.Unmarshal(body, &userToFollow)

	if &userToFollow == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(userToFollow)
		return
	}

}
