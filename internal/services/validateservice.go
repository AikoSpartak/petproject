package services

import (
	"errors"
	"regexp"

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

	if password == "" {
		return "", ErrorsPassword
	}

	if len(password) < 5 {
		return "", ErrorsPassword
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashPass), nil

}

func validateNumber(number string) error {
	if number == "" {
		return ErrorsPhone
	}

	if len(number) < 10 {
		return ErrorsPhone
	}

	regular := `^\+[1-9]\d{1,14}$`
	regularNumber := regexp.MustCompile(regular)

	if !regularNumber.MatchString(number) {
		return ErrorsPhone
	}

	return nil
}
