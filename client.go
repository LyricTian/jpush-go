package jpush

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
)

// NewClient 创建推送客户端实例
func NewClient(maxThread int, opts ...Option) *Client {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	cli := &Client{
		opts:      &o,
		queue:     NewQueue(maxThread),
		cidClient: NewCIDClient(1000, opts...),
	}

	cli.jobPool = &sync.Pool{
		New: func() interface{} {
			return newPushJob(cli.opts, cli.queue)
		},
	}
	cli.queue.Start()

	return cli
}

// Client 推送客户端
type Client struct {
	opts      *options
	queue     *Queue
	cidClient *CIDClient
	jobPool   *sync.Pool
}

func (c *Client) fillCID(ctx context.Context, payload *Payload) error {
	if payload.CID == "" {
		cid, err := c.GetPushID(ctx)
		if err != nil {
			return err
		}
		payload.CID = cid
	}
	return nil
}

// Terminate 终止客户端
func (c *Client) Terminate() {
	c.queue.Stop()
}

// GetPushID 获取推送ID
func (c *Client) GetPushID(ctx context.Context) (string, error) {
	return c.cidClient.GetPushID(ctx)
}

// GetScheduleID 获取定时ID
func (c *Client) GetScheduleID(ctx context.Context) (string, error) {
	return c.cidClient.GetScheduleID(ctx)
}

// Push 消息推送
func (c *Client) Push(ctx context.Context, payload *Payload, callback PushResultHandle) error {
	err := c.fillCID(ctx, payload)
	if err != nil {
		return err
	}

	job := c.jobPool.Get().(*pushJob)
	job.Reset(ctx, payload, callback)
	c.queue.Push(job)
	return nil
}

// PushValidate 先校验，再推送
func (c *Client) PushValidate(ctx context.Context, payload *Payload, callback PushResultHandle) error {
	err := c.fillCID(ctx, payload)
	if err != nil {
		return err
	}

	resp, err := internalRequest(ctx, c.opts, "/v3/push/validate", http.MethodPost, payload.Reader())
	if err != nil {
		return err
	}
	resp.Close()

	return c.Push(ctx, payload, callback)
}

// PushResult 推送响应结果
type PushResult struct {
	SendNO string `json:"sendno"`
	MsgID  string `json:"msg_id"`
}

func (r *PushResult) String() string {
	buf, _ := json.Marshal(r)
	return string(buf)
}

// PushResultHandle 异步响应结果
type PushResultHandle func(*PushResult, error)
