package handler

import (
	resreqstruct "arch/internal/handler/ResReqStruct"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
)

func (u *UserHandler) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {

	var req resreqstruct.UpdateProfileRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid getprofile request", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := u.UserServices.UpdateProfile(r.Context(), id, req.Name, req.Phone)
	if err != nil {
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

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response update", http.StatusInternalServerError)
		return
	}

}
