package handler

import (
	resreqstruct "arch/internal/handler/ResReqStruct"
	"arch/internal/middleware"
	"arch/internal/services"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
)

type UserHandler struct {
	UserServices *services.UserService
}

func (u *UserHandler) RegisterRoutes(r *chi.Mux) {

	r.Post("/register", u.RegisterHandler)
	r.Post("/login", u.LoginHandler)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Get("/users/{id}", u.GetUserHandler)
		r.Put("/users/{id}", u.UpdateProfileHandler)
		r.Put("/users/{id}/email", u.ChangeEmailHandler)
		r.Put("/users/{id}/password", u.HandlerChangePassword)
		r.Delete("/users/{id}", u.DeleteHandler)
	})
}

func (u *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var req resreqstruct.RegisterRequest

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
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	respnose := resreqstruct.ResponseRegister{
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
