package services

import (
	"arch/internal/models"
	"arch/internal/repository/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteUser(t *testing.T) {

	userID := uuid.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepository(ctrl)

	mockRepository.EXPECT().
		GetByUserId(gomock.Any(), userID).
		Return(&models.User{
			Uuid: userID,
		}, nil)

	mockRepository.EXPECT().
		DeleteUser(gomock.Any(), userID).
		Return(nil)

	mock := &UserService{UserRepositor: mockRepository}

	err := mock.DeleteAccountUser(t.Context(), userID)

	require.NoError(t, err)
}
