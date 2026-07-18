package handler

import (
	resreqstruct "arch/internal/handler/ResReqStruct"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (u *UserHandler) ChangeEmailHandler(w http.ResponseWriter, r *http.Request) {

	var req resreqstruct.ChangeEmailRequest

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

	user, err := u.UserServices.ChangeEmail(r.Context(), id, req.Email)
	if err != nil {
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
