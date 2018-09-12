package jpush

import (
	"bytes"
	"encoding/json"
	"io"
)

// OS 推送平台
type OS string

func (o OS) String() string {
	return string(o)
}

// 定义推送平台
const (
	Android  OS = "android"
	iOS      OS = "ios"
	WinPhone OS = "winphone"
)

// Payload 推送载荷
type Payload struct {
	Platform     *Platform     `json:"platform"`               // 推送平台
	Audience     *Audience     `json:"audience"`               // 推送目标
	Notification *Notification `json:"notification,omitempty"` // 通知
	Message      *Message      `json:"message,omitempty"`      // 自定义消息
	SmsMessage   *SmsMessage   `json:"sms_message,omitempty"`  // 短信补充
	Options      *Options      `json:"options,omitempty"`      // 可选参数
	CID          string        `json:"cid,omitempty"`          // 推送唯一标识符
}

// Reader 序列化为 JSON 流
func (p *Payload) Reader() io.Reader {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(p)
	return buf
}

// NewPlatform 创建推送平台实例
func NewPlatform() *Platform {
	return new(Platform)
}

// Platform 推送平台
type Platform struct {
	IsAll bool
	Value []string
}

// MarshalJSON 实现 JSON 接口
func (p *Platform) MarshalJSON() ([]byte, error) {
	if p.IsAll {
		return json.Marshal("all")
	}
	return json.Marshal(p.Value)
}

// All 推送到所有平台
func (p *Platform) All() *Platform {
	p.IsAll = true
	return p
}

// Add 指定特定推送平台
func (p *Platform) Add(oss ...OS) *Platform {
	for _, os := range oss {
		exists := false

		for _, v := range p.Value {
			if v == os.String() {
				exists = true
				break
			}
		}

		if !exists {
			p.Value = append(p.Value, os.String())
		}
	}
	return p
}

// NewAudience 创建推送目标实例
func NewAudience() *Audience {
	return new(Audience)
}

// Audience 推送目标
type Audience struct {
	IsAll bool
	Value map[string][]string
}

// MarshalJSON 实现 JSON 接口
func (a *Audience) MarshalJSON() ([]byte, error) {
	if a.IsAll {
		return json.Marshal("all")
	}
	return json.Marshal(a.Value)
}

// All 全部设备
func (a *Audience) All() *Audience {
	a.IsAll = true
	return a
}

// SetValue 设定推送目标
func (a *Audience) SetValue(key string, values ...string) *Audience {
	if a.Value == nil {
		a.Value = make(map[string][]string)
	}

	a.Value[key] = values
	return a
}

// SetTag 设定标签 OR
func (a *Audience) SetTag(tags ...string) *Audience {
	return a.SetValue("tag", tags...)
}

// SetTagAnd 设定标签 AND
func (a *Audience) SetTagAnd(tags ...string) *Audience {
	return a.SetValue("tag_and", tags...)
}

// SetTagNot 设定标签 NOT
func (a *Audience) SetTagNot(tags ...string) *Audience {
	return a.SetValue("tag_not", tags...)
}

// SetAlias 设定别名
func (a *Audience) SetAlias(aliases ...string) *Audience {
	return a.SetValue("alias", aliases...)
}

// SetRegistrationID 设定注册 ID
func (a *Audience) SetRegistrationID(registrationIDs ...string) *Audience {
	return a.SetValue("registration_id", registrationIDs...)
}

// SetSegment 设定用户分群 ID
func (a *Audience) SetSegment(segments ...string) *Audience {
	return a.SetValue("segment", segments...)
}

// SetAbTest 设定A/B Test ID
func (a *Audience) SetAbTest(abtests ...string) *Audience {
	return a.SetValue("abtest", abtests...)
}

// NewNotification 创建通知实例
func NewNotification() *Notification {
	return new(Notification)
}

// Notification 通知
type Notification struct {
	Alert    string                `json:"alert,omitempty"`
	Android  *AndroidNotification  `json:"android,omitempty"`
	IOS      *IOSNotification      `json:"ios,omitempty"`
	WinPhone *WinPhoneNotification `json:"winphone,omitempty"`
}

// SetAlert 设定通知内容
func (n *Notification) SetAlert(alert string) *Notification {
	n.Alert = alert
	return n
}

// SetAndroidNotification 设定 Android 平台上的通知
func (n *Notification) SetAndroidNotification(android *AndroidNotification) *Notification {
	n.Android = android
	return n
}

// SetIOSNotification 设定 iOS 平台上的通知
func (n *Notification) SetIOSNotification(ios *IOSNotification) *Notification {
	n.IOS = ios
	return n
}

// SetWinPhoneNotification 设定 Windows Phone 平台上的通知
func (n *Notification) SetWinPhoneNotification(winPhone *WinPhoneNotification) *Notification {
	n.WinPhone = winPhone
	return n
}

// NewAndroidNotification 创建 Android 平台上的通知实例
func NewAndroidNotification() *AndroidNotification {
	return new(AndroidNotification)
}

// AndroidNotification Android 平台上的通知
type AndroidNotification struct {
	Alert      string                 `json:"alert"`
	Title      string                 `json:"title,omitempty"`
	BuilderID  int                    `json:"builder_id,omitempty"`
	Priority   int                    `json:"priority,omitempty"`
	Category   string                 `json:"category,omitempty"`
	Style      int                    `json:"style,omitempty"`
	AlertType  int                    `json:"alert_type,omitempty"`
	BigText    string                 `json:"big_text,omitempty"`
	Inbox      map[string]interface{} `json:"inbox,omitempty"`
	BigPicPath string                 `json:"big_pic_path,omitempty"`
	Extras     map[string]interface{} `json:"extras,omitempty"`
}

// SetAlert 通知内容
func (n *AndroidNotification) SetAlert(alert string) *AndroidNotification {
	n.Alert = alert
	return n
}

// SetTitle 通知标题
func (n *AndroidNotification) SetTitle(title string) *AndroidNotification {
	n.Title = title
	return n
}

// SetExtras 扩展字段
func (n *AndroidNotification) SetExtras(extras map[string]interface{}) *AndroidNotification {
	n.Extras = extras
	return n
}

// NewIOSNotification 创建 iOS 平台上的通知实例
func NewIOSNotification() *IOSNotification {
	return new(IOSNotification)
}

// IOSNotification iOS 平台上 APNs 通知结构
type IOSNotification struct {
	Alert            interface{}            `json:"alert"`
	Sound            string                 `json:"sound,omitempty"`
	Badge            interface{}            `json:"badge,omitempty"`
	ContentAvailable bool                   `json:"content-available,omitempty"`
	MutableContent   bool                   `json:"mutable-content,omitempty"`
	Category         string                 `json:"category,omitempty"`
	Extras           map[string]interface{} `json:"extras,omitempty"`
}

// SetAlert 通知内容
func (n *IOSNotification) SetAlert(alert interface{}) *IOSNotification {
	n.Alert = alert
	return n
}

// SetBadge 应用角标
func (n *IOSNotification) SetBadge(badge interface{}) *IOSNotification {
	n.Badge = badge
	return n
}

// SetExtras 扩展字段
func (n *IOSNotification) SetExtras(extras map[string]interface{}) *IOSNotification {
	n.Extras = extras
	return n
}

// NewWinPhoneNotification 创建 Windows Phone 平台上的通知实例
func NewWinPhoneNotification() *WinPhoneNotification {
	return new(WinPhoneNotification)
}

// WinPhoneNotification Windows Phone 平台上的通知
type WinPhoneNotification struct {
	Alert    string                 `json:"alert"`
	Title    string                 `json:"title,omitempty"`
	OpenPage string                 `json:"_open_page,omitempty"`
	Extras   map[string]interface{} `json:"extras,omitempty"`
}

// SetAlert 通知内容
func (n *WinPhoneNotification) SetAlert(alert string) *WinPhoneNotification {
	n.Alert = alert
	return n
}

// SetTitle 通知标题
func (n *WinPhoneNotification) SetTitle(title string) *WinPhoneNotification {
	n.Title = title
	return n
}

// SetExtras 扩展字段
func (n *WinPhoneNotification) SetExtras(extras map[string]interface{}) *WinPhoneNotification {
	n.Extras = extras
	return n
}

// Message 自定义消息
type Message struct {
	Content     string                 `json:"msg_content"`
	Title       string                 `json:"title,omitempty"`
	ContentType string                 `json:"content_type,omitempty"`
	Extras      map[string]interface{} `json:"extras,omitempty"`
}

// SmsMessage 短信补充
type SmsMessage struct {
	TempPara  interface{} `json:"temp_para,omitempty"`
	TempID    int64       `json:"temp_id"`
	DelayTime int         `json:"delay_time"`
}

// NewOptions 创建可选参数实例
func NewOptions() *Options {
	return new(Options)
}

// Options 可选参数
type Options struct {
	SendNO          int    `json:"sendno,omitempty"`
	TimeLive        int    `json:"time_to_live,omitempty"`
	OverrideMsgID   int64  `json:"override_msg_id,omitempty"`
	ApnsProduction  bool   `json:"apns_production"`
	ApnsCollapseID  string `json:"apns_collapse_id,omitempty"`
	BigPushDuration int    `json:"big_push_duration,omitempty"`
}

// SetApnsProduction 设定 APNs 是否生产环境
func (o *Options) SetApnsProduction(prod bool) *Options {
	o.ApnsProduction = prod
	return o
}
