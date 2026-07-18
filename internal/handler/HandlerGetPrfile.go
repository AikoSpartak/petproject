package handler

import (
	resreqstruct "arch/internal/handler/ResReqStruct"
	"encoding/json"
	"net/http"
)

func (u *UserHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {

	var req resreqstruct.GetPrfileRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := u.UserServices.GetProfile(r.Context(), req.ID)
	if err != nil {
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

	w.Header().Set("Content-type", "application-json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
