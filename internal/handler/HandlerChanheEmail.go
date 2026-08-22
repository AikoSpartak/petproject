package handler

import (
	resreqstruct "arch/internal/handler/ResReqStruct"
	"arch/internal/middleware"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (u *UserHandler) ChangeEmailHandler(w http.ResponseWriter, r *http.Request) {

	var req resreqstruct.ChangeEmailRequest

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

	user, err := u.UserServices.ChangeEmail(r.Context(), id, req.Email)
	if err != nil {
		slog.Error("mail change error when requesting",
			"error", err.Error())
		http.Error(w, "failed to change email", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-type", "json-application")
	w.WriteHeader(http.StatusOK)

	response := resreqstruct.ChangeEmailResponse{
		Uuid:  user.Uuid,
		Email: user.Email,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}
