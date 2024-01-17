package context

import "context"

func IsAuth(ctx context.Context) bool {
	if isAuth, ok := ctx.Value(AuthContextKey).(bool); ok {
		return isAuth
	}
	return false
}

func GetPostCount(ctx context.Context) string {
	if postCount, ok := ctx.Value(PostCountContextKey).(string); ok {
		return postCount
	}
	return "???"
}

type contextKey string

var AuthContextKey contextKey = "isAuth"
var PostCountContextKey contextKey = "postCount"
