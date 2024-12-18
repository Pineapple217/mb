package context

import (
	"context"
	"strconv"
)

func IsAuth(ctx context.Context) bool {
	if isAuth, ok := ctx.Value(AuthContextKey).(bool); ok {
		return isAuth
	}
	return false
}

func GetPostCountStr(ctx context.Context) string {
	if postCount, ok := ctx.Value(PostCountContextKey).(int64); ok {
		if postCount == -1 {
			return "???"
		}
		return strconv.FormatInt(postCount, 10)
	}
	return "???"
}

func GetPostCount(ctx context.Context) int64 {
	if postCount, ok := ctx.Value(PostCountContextKey).(int64); ok {
		return postCount
	}
	return -1
}

func GetPath(ctx context.Context) string {
	if path, ok := ctx.Value(PathContextKey).(string); ok {
		return path
	}
	return ""
}

type contextKey string

var AuthContextKey contextKey = "isAuth"
var PostCountContextKey contextKey = "postCount"
var PathContextKey contextKey = "path"
