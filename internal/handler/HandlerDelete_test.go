package handler

import (
	"arch/internal/services/mocks"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestDeleteUserHandler(t *testing.T) {

	userId := uuid.New()

	ctrl := gomock.NewController(t)
	mockServices := mocks.NewMockUserSrvices(ctrl)

	mockServices.EXPECT().DeleteAccountUser(gomock.Any(), userId).Return(nil)

	handler := &UserHandler{
		UserServices: mockServices,
	}

	req := httptest.NewRequest(
		http.MethodDelete,
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

	handler.DeleteUserHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf(
			"error status delete %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

}
