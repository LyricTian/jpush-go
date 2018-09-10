package jpush

import (
	"context"
	"sync"
)

var (
	internalClient *Client
	once           sync.Once
)

// Init 初始化推送客户端
func Init(maxThread int, opts ...Option) *Client {
	once.Do(func() {
		internalClient = NewClient(maxThread, opts...)
	})
	return internalClient
}

func client() *Client {
	if internalClient == nil {
		panic("Client is not initialized")
	}
	return internalClient
}

// Terminate 终止客户端
func Terminate() {
	client().Terminate()
}

// GetPushID 获取推送ID
func GetPushID(ctx context.Context) (string, error) {
	return client().GetPushID(ctx)
}

// GetScheduleID 获取推送ID
func GetScheduleID(ctx context.Context) (string, error) {
	return client().GetScheduleID(ctx)
}

// Push 消息推送
func Push(ctx context.Context, payload *Payload, callback PushResultHandle) error {
	return client().Push(ctx, payload, callback)
}

// PushValidate 先校验，再推送
func PushValidate(ctx context.Context, payload *Payload, callback PushResultHandle) error {
	return client().PushValidate(ctx, payload, callback)
}
