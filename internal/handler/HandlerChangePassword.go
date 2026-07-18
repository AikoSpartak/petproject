package handler

import (
	resreqstruct "arch/internal/handler/ResReqStruct"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (u *UserHandler) HandlerChangePassword(w http.ResponseWriter, r *http.Request) {

	var req resreqstruct.ChangePasswordRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request email", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	err = u.UserServices.ChangePassword(r.Context(), id, req.OldPassword, req.NewPassword)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := resreqstruct.ChangePasswordResponse{
		Message: "password changed successfully",
	}

	json.NewEncoder(w).Encode(response)
}
