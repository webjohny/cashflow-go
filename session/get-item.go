package session

import (
	"github.com/gin-gonic/gin"
)

func GetItem[T comparable](ctx *gin.Context, key string) T {
	store := GetStore(ctx)

	value, ok := store.Get(key)

	if ok {
		return value.(T)
	}

	return value.(T)
}

func SetItem(ctx *gin.Context, key string, value interface{}) {
	store := GetStore(ctx)

	store.Set(key, value)
}
