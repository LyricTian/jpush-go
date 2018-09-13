package jpush

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var (
	// ErrInvalidCID 无效的cid
	ErrInvalidCID = errors.New("invalid cid")
)

// NewCIDClient 创建获取CID实例
func NewCIDClient(count int, opts ...Option) *CIDClient {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	return &CIDClient{
		opts:         &o,
		pushItem:     newCIDItem(&o, "push", count),
		scheduleItem: newCIDItem(&o, "schedule", count),
	}
}

// CIDClient 推送唯一标识符客户端
type CIDClient struct {
	opts         *options
	pushItem     *cidItem
	scheduleItem *cidItem
}

// GetPushID 获取推送ID
func (c *CIDClient) GetPushID(ctx context.Context) (string, error) {
	return c.pushItem.Get(ctx)
}

// GetScheduleID 获取定时ID
func (c *CIDClient) GetScheduleID(ctx context.Context) (string, error) {
	return c.scheduleItem.Get(ctx)
}

func newCIDItem(opts *options, typ string, count int) *cidItem {
	return &cidItem{
		opts:  opts,
		lock:  new(sync.RWMutex),
		list:  list.New(),
		typ:   typ,
		count: count,
	}
}

type cidItem struct {
	opts      *options
	lock      *sync.RWMutex
	expiredAt time.Time
	list      *list.List
	typ       string
	count     int
}

func (c *cidItem) Get(ctx context.Context) (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	ele := c.list.Front()
	// 如果没有过期，并且队列中有值，则直接返回
	if c.expiredAt.After(time.Now()) && ele != nil {
		c.list.Remove(ele)
		return ele.Value.(string), nil
	}

	params := make(url.Values)
	params.Set("type", c.typ)
	params.Set("count", strconv.Itoa(c.count))

	router := fmt.Sprintf("/v3/push/cid?%s", params.Encode())
	resp, err := pushRequest(ctx, c.opts, router, http.MethodGet, nil)
	if err != nil {
		return "", err
	}

	var result struct {
		CIDList []string `json:"cidlist"`
	}
	err = resp.JSON(&result)
	if err != nil {
		return "", err
	}

	if len(result.CIDList) > 0 {
		c.expiredAt = time.Now().Add(time.Hour * 23)
		c.list = c.list.Init()
		for _, v := range result.CIDList {
			c.list.PushBack(v)
		}
		return c.list.Front().Value.(string), nil
	}

	return "", ErrInvalidCID
}
