package handler

import (
	resreqstruct "arch/internal/handler/ResReqStruct"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (u *UserHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
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

	json.NewEncoder(w).Encode(response)
}
