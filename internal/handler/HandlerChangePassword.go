package handler

import (
	resreqstruct "arch/internal/handler/ResReqStruct"
	"arch/internal/middleware"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (u *UserHandler) HandlerChangePassword(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	var req resreqstruct.ChangePasswordRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("parsing error",
			"error", err.Error())
		http.Error(w, "invalid request email", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		slog.Error("error getting the id during the request",
			"error", err.Error())
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	userRole, ok := middleware.UserRoleFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorezed", http.StatusUnauthorized)
		return
	}

	userId, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !middleware.VerificationRole(userRole, userId, id.String()) {
		http.Error(w, "access in denied", http.StatusForbidden)
		return
	}

	err = u.UserServices.ChangePassword(r.Context(), id, req.OldPassword, req.NewPassword)
	if err != nil {
		slog.Error("error receiving the password when requesting",
			"error", err.Error())
		http.Error(w, "Failed to update password", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := resreqstruct.ChangePasswordResponse{
		Message: "password changed successfully",
	}

	slog.Info("successful password change request",
		"time", time.Since(start).Milliseconds())

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
