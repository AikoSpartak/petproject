package services

import (
	"errors"
	"regexp"

	"github.com/nyaruka/phonenumbers/v2"
	"golang.org/x/crypto/bcrypt"
)

func emailValidate(email string) error {

	if email == "" {
		return ErrorSenEmail
	}

	if len(email) < 5 {
		return ErrorSenEmail
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil

}

func hashPassword(password string) (string, error) {

	if password == "" || len(password) < 8 {
		return "", ErrorsPassword
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashPass), nil

}

func validateNumber(number string) error {
	if number == "" || len(number) < 10 {
		return ErrorsPhone
	}

	res, err := phonenumbers.Parse(number, "")
	if err != nil {
		return err
	}

	if !phonenumbers.IsValidNumber(res) {
		return ErrorsPhone
	}

	return nil
}
