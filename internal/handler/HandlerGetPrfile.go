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

func (u *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		slog.Error("error receiving the ID",
			"error", err.Error())
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	userRole, ok := middleware.UserRoleFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !middleware.VerificationRole(userRole, userID, id.String()) {
		http.Error(w, "access is denied", http.StatusForbidden)
		return
	}

	user, err := u.UserServices.GetProfile(r.Context(), id)
	if err != nil {
		slog.Error("error receiving the user during the request",
			"error", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := resreqstruct.GetPrfileResponse{
		Uuid:         user.Uuid.String(),
		Email:        user.Email,
		FullName:     user.FullName,
		PhoneNumber:  user.PhoneNumber,
		DateBirthday: user.DateBirthday,
		Role:         user.Role,
		Country:      user.Country,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
