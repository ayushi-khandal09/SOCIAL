package cache

import (
	"context"

	"github.com/ayushi-khandal09/social/internal/store"
	"github.com/go-redis/redis/v8"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
		Delete(context.Context, int64) error
	}
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rdb},
	}
}
