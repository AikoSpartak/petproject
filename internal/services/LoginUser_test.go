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

func TestLoginUser(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()

		t.Setenv("JWT_KEY", "test-secret-key-for-jwt")
		t.Setenv("SECRET_KEY", "test-secret-key-for-refresh")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		hashPassowrd, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)

		mockRepository.EXPECT().GetByEmail(gomock.Any(), "admin@mail.ru").
			Return(&models.User{
				Uuid:     userID,
				Email:    "admin@mail.ru",
				Password: string(hashPassowrd),
				Role:     "admin",
			}, nil)

		mockRepository.EXPECT().DeleteAccessToken(userID)
		mockRepository.EXPECT().DeleteRefreshToken(userID)

		mockRepository.EXPECT().SaveAccessToken(userID, gomock.Any())
		mockRepository.EXPECT().SaveRefreshToken(userID, gomock.Any())

		mock := &UserService{UserRepositor: mockRepository}

		accessToken, refreshToken, user, err := mock.LoginUser(
			t.Context(),
			"admin@mail.ru",
			"12345678")

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, userID, user.Uuid)
		require.Equal(t, "admin@mail.ru", user.Email)
		require.NotEmpty(t, accessToken)
		require.NotEmpty(t, refreshToken)
	})

	t.Run("user not found", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		mockRepository.EXPECT().
			GetByEmail(gomock.Any(), "sasa@mail.ru").
			Return(nil, nil)

		mock := &UserService{UserRepositor: mockRepository}

		access, refresh, user, err := mock.LoginUser(t.Context(), "sasa@mail.ru", "asdf")

		require.Error(t, err)
		require.Nil(t, user)
		require.Empty(t, access)
		require.Empty(t, refresh)
	})

	t.Run("error password", func(t *testing.T) {

		correctPassword := "123456789"
		noCurrectPassword := "123"

		userID := uuid.New()

		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		mockRepository.EXPECT().
			GetByEmail(gomock.Any(), "admin@mail.ru").
			Return(&models.User{
				Uuid:     userID,
				Email:    "admin@mail.ru",
				Password: string(hashPassword),
				Role:     "admin",
			}, nil)

		mock := &UserService{UserRepositor: mockRepository}

		access, refresh, user, err := mock.LoginUser(t.Context(), "admin@mail.ru", noCurrectPassword)

		require.Error(t, err)
		require.Nil(t, user)
		require.Empty(t, access)
		require.Empty(t, refresh)
	})
}
