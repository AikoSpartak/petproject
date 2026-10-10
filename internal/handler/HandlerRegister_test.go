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

func TestRegisterUserHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockServices := mocks.NewMockUserSrvices(ctrl)

	mockServices.EXPECT().RegisterUser(gomock.Any(),
		"qwerqwer@mail.ru",
		"qwerqwerqwer",
		"+79123265661",
		"1993-16-11",
		"768983",
		"Russia").Return(
		&models.User{
			Uuid:        uuid.New(),
			Email:       "qwerqwer@mail.ru",
			PhoneNumber: "+79123265661",
		},
		nil,
	)

	handler := &UserHandler{
		UserServices: mockServices,
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/register",
		strings.NewReader(`{
		"email": "qwerqwer@mail.ru",
		"password": "qwerqwerqwer",
		"phonenumber": "+79123265661",
		"datebirthday": "1993-16-11",
		"passport": "768983",
		"country": "Russia"}`),
	)

	recoder := httptest.NewRecorder()

	handler.RegisterUserHandler(recoder, req)

	if recoder.Code != http.StatusOK {
		t.Fatalf("expected status refister, got %d", recoder.Code)
	}
}
