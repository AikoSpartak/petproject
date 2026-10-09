package services

import (
	"arch/internal/repository/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRegisterUser(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		mockRepository.EXPECT().
			CreateUser(gomock.Any(), gomock.Any()).
			Return(nil)

		mock := &UserService{UserRep: mockRepository}

		user, err := mock.RegisterUser(t.Context(),
			"admin@mail.ru",
			"123456789",
			"+79995844443",
			time.Date(1993, 1, 1, 0, 0, 0, 0, time.UTC),
			"767289",
			"Russia")

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, "admin@mail.ru", user.Email)
		require.NotEqual(t, "123456789", user.Password)
		require.NotEmpty(t, user.Password)
		require.Equal(t, "+79995844443", user.PhoneNumber)
		require.Equal(t, "767289", user.Passport)
		require.Equal(t, "Russia", user.Country)
		require.Equal(t, time.Date(1993, 1, 1, 0, 0, 0, 0, time.UTC), user.DateBirthday)
	})

	t.Run("invalid email", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		mock := &UserService{UserRep: mockRepository}

		user, err := mock.RegisterUser(t.Context(),
			"admin@",
			"123456789",
			"+79995844443",
			time.Date(1993, 1, 1, 0, 0, 0, 0, time.UTC),
			"786564",
			"Russia")
		require.Error(t, err)
		require.Nil(t, user)
	})

	t.Run("invalid password", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		mock := &UserService{UserRep: mockRepository}

		user, err := mock.RegisterUser(t.Context(),
			"admin@mail.ru",
			"12345",
			"+79995844443",
			time.Date(1993, 1, 1, 0, 0, 0, 0, time.UTC),
			"786564",
			"Russia")
		require.Error(t, err)
		require.Nil(t, user)

	})

	t.Run("invalid humber", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockUserRepository(ctrl)

		mock := &UserService{UserRep: mockRepository}

		user, err := mock.RegisterUser(t.Context(),
			"admin@mail.ru",
			"123456789",
			"123",
			time.Date(1993, 1, 1, 0, 0, 0, 0, time.UTC),
			"786564",
			"Russia")
		require.Error(t, err)
		require.Nil(t, user)
	})

}
