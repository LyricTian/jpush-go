package jpush

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/LyricTian/queue"
)

func newPushJob(opts *options, queue queue.Queuer) *pushJob {
	return &pushJob{
		opts:  opts,
		queue: queue,
	}
}

type pushJob struct {
	opts     *options
	queue    queue.Queuer
	payload  *Payload
	ctx      context.Context
	callback PushResultHandle
}

func (j *pushJob) Reset(ctx context.Context, payload *Payload, callback PushResultHandle) {
	j.payload = payload
	j.ctx = ctx
	j.callback = callback
}

func (j *pushJob) Job() {
	resp, err := pushRequest(j.ctx, j.opts, "/v3/push", http.MethodPost, j.payload.Reader())
	if err != nil {
		if e, ok := err.(*Error); ok {
			// 如果当前推送频次超出限制，则将任务重新放入队列，并休眠等待
			if e.StatusCode == 429 && e.HeaderItem.XRateLimitReset > 0 {
				j.queue.Push(j)
				time.Sleep(time.Second * time.Duration(e.HeaderItem.XRateLimitReset))
			} else {
				j.callback(nil, err)
			}
		} else {
			if strings.Contains(err.Error(), "connection refused") {
				j.queue.Push(j)
				return
			}
			j.callback(nil, err)
		}
		return
	}

	var result PushResult
	err = resp.JSON(&result)
	if err != nil {
		j.callback(nil, err)
		return
	}
	j.callback(&result, nil)
}
