package handler

import (
	"net/http"

	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		DB: db,
	}
}

// создали ручку для проверки подключения к бд
func (h *Handler) HandlerApp(w http.ResponseWriter, r *http.Request) {

	var result int

	if err := h.DB.Raw("SELECT 1").Scan(&result).Error; err != nil {
		response := ResponseHandler("connection failed with an error", err.Error())
		writeJson(w, http.StatusInternalServerError, response)
		return
	}

	response := ResponseHandler("connection was successful", "")
	writeJson(w, http.StatusOK, response)

}
