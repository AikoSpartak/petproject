package handler

import (
	"arch/internal/models"
	"arch/internal/services/mocks"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestUserGetProfileHandler(t *testing.T) {

	userId := uuid.New()

	ctrl := gomock.NewController(t)
	mockServices := mocks.NewMockUserSrvices(ctrl)

	mockServices.EXPECT().GetUser(gomock.Any(), userId).Return(
		&models.User{
			Uuid:        userId,
			FullName:    "Gena",
			Email:       "qwerqwer@mail.ru",
			PhoneNumber: "+79123265661",
			Role:        "user",
			Country:     "Russia",
		},
		nil,
	)

	handler := &UserHandler{
		UserServices: mockServices,
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/"+userId.String(),
		nil,
	)

	ctx := context.WithValue(
		req.Context(),
		"user_id",
		userId.String(),
	)

	ctx = context.WithValue(
		ctx,
		"role",
		"user",
	)

	req = req.WithContext(ctx)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", userId.String())

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	recorder := httptest.NewRecorder()

	handler.GetUserHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf(
			"error status profile %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

}
