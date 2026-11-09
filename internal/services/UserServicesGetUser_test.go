package services

import (
	"arch/internal/models"
	"arch/internal/repository/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUser(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		mockRepository.EXPECT().
			GetByUserId(gomock.Any(), userID).
			Return(&models.User{
				Uuid:        userID,
				Email:       "admin@mail.ru",
				PhoneNumber: "+79995844443",
			}, nil)

		mock := &UserService{UserRepositor: mockRepository}

		user, err := mock.GetUser(t.Context(), userID)

		require.NoError(t, err)
		require.NotNil(t, user)
	})

	t.Run("user not found", func(t *testing.T) {

		userID := uuid.New()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		mockRepository.EXPECT().
			GetByUserId(gomock.Any(), userID).
			Return(nil, nil)

		mock := &UserService{UserRepositor: mockRepository}

		user, err := mock.GetUser(t.Context(), userID)

		require.Error(t, err)
		require.Nil(t, user)

	})
}
