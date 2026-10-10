package handler

import (
	"arch/internal/handler/RequestStruct"
	"arch/internal/handler/ResponseStruct"
	"encoding/json"
	"log/slog"
	"net/http"
)

func (u *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var req RequestStruct.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("error during authorization",
			"error", err.Error())
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, user, err := u.UserServices.LoginUser(
		r.Context(),
		req.Login,
		req.Password,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := ResponseStruct.LoginResponse{
		Uuid:         user.Uuid.String(),
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
