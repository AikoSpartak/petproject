package services

import (
	"arch/internal/models"
	"arch/internal/repository/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestChangeEmailUser(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()

		newEmail := "admin@mail.ru"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		mockRepository.EXPECT().
			GetByEmail(gomock.Any(), newEmail).
			Return(&models.User{
				Uuid:  userID,
				Email: "admin1@mail.ru",
			}, nil)

		mockRepository.EXPECT().
			GetByUserId(gomock.Any(), userID).
			Return(&models.User{
				Uuid:  userID,
				Email: "admin@mail.ru",
			}, nil)

		mockRepository.EXPECT().
			UpdateUserEmail(gomock.Any(), gomock.Any()).
			Return(nil)

		mock := &UserService{UserRepositor: mockRepository}

		user, err := mock.ChangeEmailUser(t.Context(), userID, newEmail)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, newEmail, user.Email)
	})

	t.Run("invalid email", func(t *testing.T) {

		userID := uuid.New()
		newEmail := "admsdfa.ru"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		mock := &UserService{UserRepositor: mockRepository}

		user, err := mock.ChangeEmailUser(t.Context(), userID, newEmail)

		require.Error(t, err)
		require.Nil(t, user)

	})

	t.Run("invalid id", func(t *testing.T) {

		oldID := uuid.New()
		newID := uuid.New()
		newEmail := "admin@mail.ru"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		mockRepository.EXPECT().
			GetByEmail(gomock.Any(), newEmail).
			Return(&models.User{
				Uuid:  oldID,
				Email: newEmail,
			}, nil)

		mock := &UserService{UserRepositor: mockRepository}

		user, err := mock.ChangeEmailUser(t.Context(), newID, newEmail)
		require.Error(t, err)
		require.Nil(t, user)

	})

}
