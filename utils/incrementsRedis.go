package utils

import (
	"backend/config"
	"context"
)

func IncrementClicksRedis(slug string) {
	ctx := context.Background()
	key := "link:" + slug + ":clicks"
	config.Rdb.Incr(ctx, key)
}