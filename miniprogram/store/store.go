package store

import (
	"github.com/silenceper/wechat/v2/miniprogram/context"
)

// UpdatableMessage 动态消息
type Store struct {
	*context.Context
}

// NewStore 实例化
func NewStore(ctx *context.Context) *Store {
	return &Store{
		Context: ctx,
	}
}
