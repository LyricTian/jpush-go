# JPush Golang Client

[![Build][Build-Status-Image]][Build-Status-Url] [![Codecov][codecov-image]][codecov-url] [![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

## 快速开始

### 下载安装

```bash
$ go get -u -v github.com/LyricTian/jpush-go
```

### 创建文件 `push.go`

```go
package main

import (
	"context"
	"fmt"

	"github.com/LyricTian/jpush-go"
)

func main() {
	jpush.Init(2,
		jpush.SetAppKey("b1ccd0dd04ec36b66c75e99f"),
		jpush.SetMasterSecret("ed431429270144d3ed53555b"),
	)

	defer jpush.Terminate()

	payload := &jpush.Payload{
		Platform: jpush.NewPlatform().All(),
		Audience: jpush.NewAudience().All(),
		Notification: &jpush.Notification{
			Alert: "推送通知测试",
		},
		Options: &jpush.Options{
			SendNO: 1,
		},
	}
	err := jpush.Push(context.Background(), payload, func(result *jpush.PushResult, err error) {
		if err != nil {
			panic(err)
		}
		fmt.Println(result.String())
	})
	if err != nil {
		panic(err)
	}
}

```

### 编译运行

```bash
$ go build push.go
$ ./push
```

> 输出结果
```json
{"sendno":"1","msg_id":"3866336947"}
```

## 特性

- 支持异步推送队列
- 自动处理频率限制
- 自动维护cid池

## MIT License

    Copyright (c) 2018 Lyric

[Build-Status-Url]: https://travis-ci.org/LyricTian/jpush-go
[Build-Status-Image]: https://travis-ci.org/LyricTian/jpush-go.svg?branch=master
[codecov-url]: https://codecov.io/gh/LyricTian/jpush-go
[codecov-image]: https://codecov.io/gh/LyricTian/jpush-go/branch/master/graph/badge.svg
[reportcard-url]: https://goreportcard.com/report/github.com/LyricTian/jpush-go
[reportcard-image]: https://goreportcard.com/badge/github.com/LyricTian/jpush-go
[godoc-url]: https://godoc.org/github.com/LyricTian/jpush-go
[godoc-image]: https://godoc.org/github.com/LyricTian/jpush-go?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg
