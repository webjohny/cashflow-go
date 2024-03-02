package session

import (
	"github.com/gin-gonic/gin"
)

func GetItem[T comparable](ctx *gin.Context, key string) *T {
	store := GetStore(ctx)

	value, ok := store.Get(key)

	if ok && value != nil {
		val := value.(T)
		return &val
	}

	return nil
}

func SetItem(ctx *gin.Context, key string, value interface{}) {
	store := GetStore(ctx)

	store.Set(key, value)
}

func DeleteItem(ctx *gin.Context, key string) {
	store := GetStore(ctx)

	store.Delete(key)
}
