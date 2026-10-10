package middleware

import "context"

func GetUserRole(ctx context.Context) (userId, role string, ok bool) {

	userId, ok = UserIDFromContext(ctx)
	if !ok {
		return "", "", false
	}

	role, ok = UserRoleFromContext(ctx)
	if !ok {
		return "", "", false
	}

	return userId, role, true

}
