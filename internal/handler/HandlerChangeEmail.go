package handler

import (
	"arch/internal/handler/RequestStruct"
	"arch/internal/handler/ResponseStruct"
	"arch/internal/middleware"
	"arch/internal/models"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (u *UserHandler) ChangeUserEmailHandler(w http.ResponseWriter, r *http.Request) {

	var req RequestStruct.ChangeEmailRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("mail parsing error when requesting",
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

	userId, userRole, ok := middleware.GetUserRole(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !middleware.VerificationRole(models.Role(userRole), userId, id.String()) {
		http.Error(w, "access in denied", http.StatusForbidden)
		return
	}

	user, err := u.UserServices.ChangeEmailUser(r.Context(), id, req.Email)
	if err != nil {
		slog.Error("mail change error when requesting",
			"error", err.Error())
		http.Error(w, "failed to change email", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-type", "json-application")
	w.WriteHeader(http.StatusOK)

	response := ResponseStruct.ChangeEmailResponse{
		Uuid:  user.Uuid,
		Email: user.Email,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}
