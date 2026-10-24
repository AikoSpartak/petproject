package handler

import (
	"arch/internal/models"
	"arch/internal/services/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestRegisterUserHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockServices := mocks.NewMockUserSrvices(ctrl)

	mockServices.EXPECT().RegisterUser(gomock.Any(),
		"qwerqwer@mail.ru",
		"qwerqwerqwer!A",
		"+79123265661",
		time.Date(1993, 1, 1, 0, 0, 0, 0, time.UTC),
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
		"password": "qwerqwerqwer!A",
		"phonenumber": "+79123265661",
		"datebirthday": "1993-01-01T00:00:00Z",
		"passport": "768983",
		"country": "Russia"}`),
	)

	recoder := httptest.NewRecorder()

	handler.RegisterUserHandler(recoder, req)

	if recoder.Code != http.StatusOK {
		t.Fatalf("expected status refister, got %d", recoder.Code)
	}
}
