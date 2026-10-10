package middleware

import "arch/internal/models"

func VerificationRole(role models.Role, jwtId, userId string) bool {
	return role.HasAccess(jwtId, userId)
}
