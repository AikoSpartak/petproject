package handler

import (
	"arch/internal/models"
	"arch/internal/services/mocks"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestChangeEmailUserHandler(t *testing.T) {

	userId := uuid.New()

	ctrl := gomock.NewController(t)
	mockServices := mocks.NewMockUserSrvices(ctrl)

	mockServices.EXPECT().ChangeEmailUser(
		gomock.Any(),
		userId,
		"asdfasdf@mail.ru").
		Return(&models.User{
			Email: "asdfasdf@mail.ru",
		}, nil)

	handler := &UserHandler{
		UserServices: mockServices,
	}

	req := httptest.NewRequest(
		http.MethodPut,
		"/users/"+userId.String()+"/email",
		strings.NewReader(`{"email": "asdfasdf@mail.ru"}`),
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

	handler.ChangeUserEmailHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("error status changeemail %d %d", recorder.Code, http.StatusOK)
	}

}
