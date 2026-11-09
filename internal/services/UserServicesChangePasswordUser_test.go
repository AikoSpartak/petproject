package services

import (
	"arch/internal/models"
	"arch/internal/repository/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestChangePasswordUser(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		oldPassword := "123456789!A"
		newPassword := "987654321!A"

		t.Setenv("JWT_KEY", "test-secret-key-for-jwt")
		t.Setenv("SECRET_KEY", "test-secret-key-for-refresh")

		hashPasOld, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		mockRepository.EXPECT().
			GetByUserId(gomock.Any(), userID).
			Return(&models.User{
				Uuid:     userID,
				Password: string(hashPasOld),
				Role:     "admin",
			}, nil)

		mockRepository.EXPECT().
			UpdateUserPassword(gomock.Any(), gomock.Any()).
			Return(nil)

		mockRepository.EXPECT().
			DeleteAccessToken(gomock.Any())

		mock := &UserService{UserRepositor: mockRepository}

		err := mock.ChangePasswordUser(t.Context(), userID, oldPassword, newPassword)

		require.NoError(t, err)

	})

	t.Run("invalid new password", func(t *testing.T) {

		userID := uuid.New()
		oldPassword := "123456789"
		newPassword := "123"

		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		mockRepository.EXPECT().
			GetByUserId(gomock.Any(), userID).
			Return(&models.User{
				Uuid:     userID,
				Password: string(hashPassword)},
				nil)

		mock := &UserService{UserRepositor: mockRepository}

		err := mock.ChangePasswordUser(t.Context(), userID, oldPassword, newPassword)

		require.Error(t, err)

	})

	t.Run("user not found", func(t *testing.T) {

		userID := uuid.New()
		oldPassword := "123456789"
		newPassword := "987654321"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		mockRepository.EXPECT().
			GetByUserId(gomock.Any(), userID).
			Return(nil, nil)

		mock := &UserService{UserRepositor: mockRepository}

		err := mock.ChangePasswordUser(t.Context(), userID, oldPassword, newPassword)

		require.Error(t, err)

	})

}
