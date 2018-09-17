package jpush

var defaultOptions = options{
	host:     "https://api.jpush.cn",
	cidCount: 1000,
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

// SetCIDCount 设定每次获取CID的数量
func SetCIDCount(count int) Option {
	return func(o *options) {
		o.cidCount = count
	}
}

type options struct {
	host         string
	appKey       string
	masterSecret string
	cidCount     int
}
