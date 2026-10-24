package handler

import (
	"arch/internal/models"
	"arch/internal/services/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestLoginUserHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockservices := mocks.NewMockUserSrvices(ctrl)

	userID := uuid.New()

	mockservices.EXPECT().LoginUser(
		gomock.Any(),
		"qwerqwer@mail.ru",
		"qwerqwerqwer").Return(
		"testaccesstoken",
		"testrefreshtoken",
		&models.User{
			Uuid:  userID,
			Email: "qwerqwer@mail.ru",
		},
		nil,
	)

	handler := &UserHandler{
		UserServices: mockservices,
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(`{
		"email": "qwerqwer@mail.ru",
		"password": "qwerqwerqwer"}`),
	)

	recorder := httptest.NewRecorder()

	handler.LoginHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

}
