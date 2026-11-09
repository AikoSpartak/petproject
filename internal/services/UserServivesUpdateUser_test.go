package services

import (
	"arch/internal/models"
	"arch/internal/repository/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdateUser(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		mockRepository.EXPECT().
			GetByUserId(gomock.Any(), userID).
			Return(&models.User{
				Uuid:        userID,
				FullName:    "aleksandr",
				PhoneNumber: "+799995844443",
			}, nil).
			Times(2)

		mockRepository.EXPECT().
			UpdateUser(gomock.Any(), userID, gomock.Any()).
			Return(nil)

		srv := &UserService{UserRepositor: mockRepository}

		user, err := srv.UpdateUser(t.Context(), userID, map[string]interface{}{
			"full_name":    "aleksandr",
			"phone_number": "+799995844443",
		})

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, "aleksandr", user.FullName)
		require.Equal(t, "+799995844443", user.PhoneNumber)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := uuid.New()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		mockRepository.EXPECT().
			GetByUserId(gomock.Any(), userID).
			Return(nil, nil)

		srv := &UserService{UserRepositor: mockRepository}

		updateMap := map[string]interface{}{
			"full_name": "gena",
		}

		user, err := srv.UpdateUser(t.Context(), userID, updateMap)

		require.Error(t, err)
		require.Nil(t, user)
	})

}
