package handler

import (
	resreqstruct "arch/internal/handler/ResReqStruct"
	"arch/internal/middleware"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
)

func (u *UserHandler) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {

	var req resreqstruct.UpdateProfileRequest
	start := time.Now()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("parsing error when updating",
			"error", err.Error())
		http.Error(w, "invalid getprofile request", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		slog.Error("error receiving the ID during the update",
			"erorr", err.Error())
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	userRole, ok := middleware.UserRoleFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userId, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !middleware.VerificationRole(userRole, userId, id.String()) {
		http.Error(w, "access is denied", http.StatusForbidden)
		return
	}

	user, err := u.UserServices.UpdateProfile(r.Context(), id, req.Name, req.Phone)
	if err != nil {
		slog.Error("error updating the user during the request",
			"error", err.Error())
		http.Error(w, "failed to update user", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := resreqstruct.UpdateProfileResponse{
		Uuid:        user.Uuid,
		FullName:    user.FullName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Country:     user.Country,
	}

	slog.Info("the user has been updated",
		"time", time.Since(start).Milliseconds())

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response update", http.StatusInternalServerError)
		return
	}

}
