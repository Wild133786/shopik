package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var reqUser CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&reqUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	respUser := CreateUserResponse{
		ID:    int64(rand.Int()),
		Email: reqUser.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respUser)
}
