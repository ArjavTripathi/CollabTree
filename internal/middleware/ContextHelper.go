package middleware

import (
	"context"
	"errors"
)

func GetUserIdFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(userIDKey).(int64)
	if !ok {
		return 0, errors.New("user id not found")
	}
	return userID, nil
}
