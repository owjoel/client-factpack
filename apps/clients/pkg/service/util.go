package service

import (
	"context"
)

func GetUsername(ctx context.Context) string {
	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return "Unknown"
	}
	return username
}
