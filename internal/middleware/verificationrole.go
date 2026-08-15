package middleware

func VerificationRole(role, jwtId, userId string) bool {

	if role == "admin" {
		return true
	}

	if role == "user" && jwtId == userId {
		return true
	}

	return false
}
