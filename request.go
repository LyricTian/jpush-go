package jpush

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
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
func internalRequest(ctx context.Context, opts *options, router, method string, body io.Reader) (Responser, error) {
	urlStr := RequestURL(opts.host, router)
	res, err := Request(ctx, urlStr, method, body, func(req *http.Request) (*http.Request, error) {
		req.SetBasicAuth(opts.appKey, opts.masterSecret)
		return req, nil
	})
	if err != nil {
		return nil, err
	} else if code := res.Response().StatusCode; code != 200 {
		e := &Error{
			StatusCode: code,
		}

		var result struct {
			Error *ErrorItem `json:"error"`
		}
		if err := res.JSON(&result); err != nil {
			return nil, err
		}
		e.ErrorItem = result.Error

		if code == 429 {
			header := new(HeaderItem)
			header.XRateLimitQuota, _ = strconv.Atoi(res.Response().Header.Get("X-Rate-Limit-Quota"))
			header.XRateLimitRemaining, _ = strconv.Atoi(res.Response().Header.Get("X-Rate-Limit-Remaining"))
			header.XRateLimitReset, _ = strconv.Atoi(res.Response().Header.Get("X-Rate-Limit-Reset"))
			e.HeaderItem = header
		}

		return nil, e
	}

	return res, nil
}

// RequestOption 自定义处理请求
type RequestOption func(*http.Request) (*http.Request, error)

// RequestURL get request url
func RequestURL(base, router string) string {
	var buf bytes.Buffer
	if l := len(base); l > 0 {
		if base[l-1] == '/' {
			base = base[:l-1]
		}
		buf.WriteString(base)

		if rl := len(router); rl > 0 {
			if router[0] != '/' {
				buf.WriteByte('/')
			}
		}
	}
	buf.WriteString(router)
	return buf.String()
}

// Request HTTP请求
func Request(ctx context.Context, urlStr, method string, body io.Reader, options ...RequestOption) (Responser, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	if len(options) > 0 {
		req, err = options[0](req)
		if err != nil {
			return nil, err
		}
	}

	var res Responser
	err = request(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		res = newResponse(resp)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func request(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{Transport: tr}
	c := make(chan error, 1)
	go func() { c <- f(client.Do(req)) }()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c
		return ctx.Err()
	case err := <-c:
		return err
	}
}

// Responser HTTP response interface
type Responser interface {
	String() (string, error)
	Bytes() ([]byte, error)
	JSON(v interface{}) error
	Response() *http.Response
	Close()
}

func newResponse(resp *http.Response) *response {
	return &response{resp}
}

type response struct {
	resp *http.Response
}

func (r *response) Response() *http.Response {
	return r.resp
}

func (r *response) String() (string, error) {
	b, err := r.Bytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (r *response) Bytes() ([]byte, error) {
	defer r.resp.Body.Close()

	buf, err := ioutil.ReadAll(r.resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (r *response) JSON(v interface{}) error {
	defer r.resp.Body.Close()

	return json.NewDecoder(r.resp.Body).Decode(v)
}

func (r *response) Close() {
	if !r.resp.Close {
		r.resp.Body.Close()
	}
}
