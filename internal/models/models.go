package models

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func (r Role) HasAccess(jwtId, userId string) bool {

	switch r {
	case RoleAdmin:
		return true
	case RoleUser:
		return jwtId == userId
	default:
		return false
	}
}
