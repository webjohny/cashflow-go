package objects

import (
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
)

type Session struct{}

func (s *Session) GetStore(ctx *gin.Context) session.Store {
	return ctx.MustGet("sessionStore").(session.Store)
}

func (s *Session) GetItem(ctx *gin.Context, key string) *string {
	store := s.GetStore(ctx)

	value, ok := store.Get(key)

	if ok && value != nil {
		val := value.(string)
		return &val
	}

	return nil
}

func (s *Session) SetItem(ctx *gin.Context, key string, value interface{}) {
	store := s.GetStore(ctx)

	store.Set(key, value)
}

func (s *Session) DeleteItem(ctx *gin.Context, key string) {
	store := s.GetStore(ctx)

	store.Delete(key)
}
