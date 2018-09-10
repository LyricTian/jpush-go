package jpush

var defaultOptions = options{
	host: "https://api.jpush.cn",
}

// Option 配置项
type Option func(o *options)

// SetHost 设定请求地址
func SetHost(host string) Option {
	return func(o *options) {
		o.host = host
	}
}

// SetAppKey 设定 Appkey
func SetAppKey(appKey string) Option {
	return func(o *options) {
		o.appKey = appKey
	}
}

// SetMasterSecret 设定 MasterSecret
func SetMasterSecret(masterSecret string) Option {
	return func(o *options) {
		o.masterSecret = masterSecret
	}
}

type options struct {
	host         string
	appKey       string
	masterSecret string
}
