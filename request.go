package jpush

import (
	"context"
	"encoding/json"
	"io"
	"strconv"

	"github.com/LyricTian/req"
)

// Error 错误
type Error struct {
	StatusCode int         `json:"status_code"`
	ErrorItem  *ErrorItem  `json:"error,omitempty"`
	HeaderItem *HeaderItem `json:"header,omitempty"`
}

func (e *Error) Error() string {
	buf, _ := json.MarshalIndent(e, "", " ")
	return string(buf)
}

// ErrorItem 错误项
type ErrorItem struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// HeaderItem 响应头
type HeaderItem struct {
	XRateLimitQuota     int `json:"X-Rate-Limit-Quota"`
	XRateLimitRemaining int `json:"X-Rate-Limit-Remaining"`
	XRateLimitReset     int `json:"X-Rate-Limit-Reset"`
}

// jpush request
func pushRequest(ctx context.Context, opts *options, router, method string, body io.Reader) (req.Responser, error) {
	urlStr := req.RequestURL(opts.host, router)
	resp, err := req.Do(ctx, urlStr, method, body, req.SetBasicAuth(opts.appKey, opts.masterSecret))
	if err != nil {
		return nil, err
	} else if code := resp.StatusCode(); code != 200 {
		var result struct {
			Error *ErrorItem `json:"error"`
		}

		if err := resp.JSON(&result); err != nil {
			return nil, err
		}

		e := &Error{
			StatusCode: code,
			ErrorItem:  result.Error,
		}

		if code == 429 {
			header := new(HeaderItem)
			header.XRateLimitQuota, _ = strconv.Atoi(resp.Response().Header.Get("X-Rate-Limit-Quota"))
			header.XRateLimitRemaining, _ = strconv.Atoi(resp.Response().Header.Get("X-Rate-Limit-Remaining"))
			header.XRateLimitReset, _ = strconv.Atoi(resp.Response().Header.Get("X-Rate-Limit-Reset"))
			e.HeaderItem = header
		}

		return nil, e
	}

	return resp, nil
}
