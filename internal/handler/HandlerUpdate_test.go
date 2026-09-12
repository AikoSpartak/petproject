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

func TestUpdateUserHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockServices := mocks.NewMockUserSrvices(ctrl)

	userID := uuid.New()

	mockServices.EXPECT().
		UpdateUser(
			gomock.Any(),
			userID,
			"Alexander",
			"+79123265661",
		).
		Return(
			&models.User{
				Uuid:        userID,
				FullName:    "Alexander",
				Email:       "qwerqwer@mail.ru",
				PhoneNumber: "+79123265661",
				Country:     "Russia",
			},
			nil,
		)

	handler := &UserHandler{
		UserServices: mockServices,
	}

	req := httptest.NewRequest(
		http.MethodPut,
		"/users/"+userID.String(),
		strings.NewReader(`{
			"name": "Alexander",
			"phone": "+79123265661"
		}`),
	)

	ctx := context.WithValue(
		req.Context(),
		"user_id",
		userID.String(),
	)

	ctx = context.WithValue(
		ctx,
		"role",
		"user",
	)

	req = req.WithContext(ctx)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", userID.String())

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	recorder := httptest.NewRecorder()

	handler.UpdateUsereHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}
}
