package handler

import (
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

func TestChangePasswordUserHandler(t *testing.T) {

	userId := uuid.New()

	ctrl := gomock.NewController(t)
	mockServices := mocks.NewMockUserSrvices(ctrl)

	mockServices.EXPECT().
		ChangePasswordUser(
			gomock.Any(),
			userId,
			"testtest",
			"test1test1").
		Return(nil)

	handler := &UserHandler{
		UserServices: mockServices,
	}

	req := httptest.NewRequest(
		http.MethodPut,
		"/users"+userId.String()+"/password",
		strings.NewReader(`{"oldpassword": "testtest",
		"newpassword": "test1test1"}`),
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

	handler.ChangeUserPasswordHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("error status changepassword %d %d", recorder.Code, http.StatusOK)
	}

}
