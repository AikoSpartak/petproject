package handler

import (
	"arch/internal/handler/ResponseStruct"
	"arch/internal/middleware"
	"arch/internal/models"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (u *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		slog.Error("error requesting an ID when deleting a user",
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
		http.Error(w, "access is denied", http.StatusForbidden)
		return
	}

	err = u.UserServices.DeleteAccountUser(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to delete account", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := ResponseStruct.DeleteResponse{
		Message: "account deleted successfully",
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
