package handler

import (
	"arch/internal/handler/RequestStruct"
	"arch/internal/handler/ResponseStruct"
	"arch/internal/middleware"
	"arch/internal/services"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
)

type UserHandler struct {
	UserServices services.UserSrvices
}

func (u *UserHandler) RegisterRoutes(r *chi.Mux) {

	r.Post("/register", u.RegisterUserHandler)
	r.Post("/login", u.LoginHandler)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Get("/users", u.GetListsUsersHandler)
		r.Get("/users/{id}", u.GetUserHandler)
		r.Put("/users/{id}/Update", u.UpdateUsereHandler)
		r.Put("/users/{id}/Email", u.ChangeUserEmailHandler)
		r.Put("/users/{id}/ChangePassword", u.ChangeUserPasswordHandler)
		r.Delete("/users/{id}", u.DeleteUserHandler)
	})
}

func (u *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {

	var req RequestStruct.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("parsing error during registration",
			"error", err.Error())
		http.Error(w, "invalid register", http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v\n", req)

	user, err := u.UserServices.RegisterUser(r.Context(),
		req.Email,
		req.Password,
		req.PhoneNumber,
		req.DateBirthday,
		req.Passport,
		req.Country)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respnose := ResponseStruct.ResponseRegister{
		Uuid:        user.Uuid,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(respnose)
	if err != nil {
		http.Error(w, "failed to encode response register", http.StatusInternalServerError)
		return
	}
}
