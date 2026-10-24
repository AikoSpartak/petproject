package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (u *UserHandler) GetListsUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	maxLimit := "50"

	limit := query.Get("limit")
	if limit == "" {
		limit = maxLimit
	}

	role := query.Get("role")
	country := query.Get("country")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "limit zero", http.StatusBadRequest)
		return
	}

	if limitInt < 1 || limitInt > 50 {
		http.Error(w, "incurrect limit", http.StatusBadRequest)
		return
	}

	listsUser, err := u.UserServices.GetListsUsers(
		r.Context(),
		false,
		limitInt,
		0,
		role,
		country,
		"")
	if err != nil {
		http.Error(w, "zero 1", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(listsUser)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}
