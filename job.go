package jpush

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/LyricTian/queue"
)

func newPushJob(opts *options, queue queue.Queuer, cidClient *CIDClient) *pushJob {
	return &pushJob{
		opts:      opts,
		queue:     queue,
		cidClient: cidClient,
	}
}

type pushJob struct {
	opts      *options
	queue     queue.Queuer
	cidClient *CIDClient
	payload   *Payload
	ctx       context.Context
	callback  PushResultHandle
}

func (j *pushJob) Reset(ctx context.Context, payload *Payload, callback PushResultHandle) {
	j.payload = payload
	j.ctx = ctx
	j.callback = callback
}

func (j *pushJob) handleError(err error) {
	if err == nil {
		return
	}

	if e, ok := err.(*Error); ok {
		if e.StatusCode == 429 || e.StatusCode == 404 {
			j.queue.Push(j)
			// 如果当前推送频次超出限制，则将任务重新放入队列，并休眠等待
			if e.HeaderItem != nil && e.HeaderItem.XRateLimitReset > 0 {
				time.Sleep(time.Second * time.Duration(e.HeaderItem.XRateLimitReset))
			}
			return
		}
		j.callback(j.ctx, nil, err)
	} else {
		if v := err.Error(); strings.Contains(v, "connection refused") {
			j.queue.Push(j)
			return
		}
		j.callback(j.ctx, nil, err)
	}
}

func (j *pushJob) Job() {
	if j.payload.CID == "" {
		cid, err := j.cidClient.GetPushID(j.ctx)
		if err != nil {
			j.handleError(err)
			return
		}
		j.payload.CID = cid
	}

	resp, err := pushRequest(j.ctx, j.opts, "/v3/push", http.MethodPost, j.payload.Reader())
	if err != nil {
		j.handleError(err)
		return
	}

	result := new(PushResult)
	err = resp.JSON(result)
	if err != nil {
		j.callback(j.ctx, nil, err)
		return
	}

	j.callback(j.ctx, result, nil)
}
