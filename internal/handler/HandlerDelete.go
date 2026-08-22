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

func (u *UserHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		slog.Error("error requesting an ID when deleting a user",
			"error", err.Error())
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
	}

	err = u.UserServices.DeleteAccount(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to delete account", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := resreqstruct.DeleteResponse{
		Message: "account deleted successfully",
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
