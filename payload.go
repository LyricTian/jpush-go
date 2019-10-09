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
	IOS      OS = "ios"
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

func (p *Payload) String() string {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(p)
	return buf.String()
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

// SetBuilderID 通知栏样式 ID
func (n *AndroidNotification) SetBuilderID(builderID int) *AndroidNotification {
	n.BuilderID = builderID
	return n
}

// SetPriority 通知栏展示优先级
func (n *AndroidNotification) SetPriority(priority int) *AndroidNotification {
	n.Priority = priority
	return n
}

// SetCategory 通知栏条目过滤或排序
func (n *AndroidNotification) SetCategory(category string) *AndroidNotification {
	n.Category = category
	return n
}

// SetStyle 通知栏样式类型
func (n *AndroidNotification) SetStyle(style int) *AndroidNotification {
	n.Style = style
	return n
}

// SetAlertType 通知提醒方式
func (n *AndroidNotification) SetAlertType(alertType int) *AndroidNotification {
	n.AlertType = alertType
	return n
}

// SetBigText 大文本通知栏样式
func (n *AndroidNotification) SetBigText(bigText string) *AndroidNotification {
	n.BigText = bigText
	return n
}

// SetInbox 文本条目通知栏样式
func (n *AndroidNotification) SetInbox(inbox map[string]interface{}) *AndroidNotification {
	n.Inbox = inbox
	return n
}

// SetBigPicPath 大图片通知栏样式
func (n *AndroidNotification) SetBigPicPath(bigPicPath string) *AndroidNotification {
	n.BigPicPath = bigPicPath
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

// SetSound 通知提示声音
func (n *IOSNotification) SetSound(sound string) *IOSNotification {
	n.Sound = sound
	return n
}

// SetBadge 应用角标
func (n *IOSNotification) SetBadge(badge interface{}) *IOSNotification {
	n.Badge = badge
	return n
}

// SetContentAvailable 推送唤醒
func (n *IOSNotification) SetContentAvailable(contentAvailable bool) *IOSNotification {
	n.ContentAvailable = contentAvailable
	return n
}

// SetMutableContent 通知扩展
func (n *IOSNotification) SetMutableContent(mutableContent bool) *IOSNotification {
	n.MutableContent = mutableContent
	return n
}

// SetCategory 通知栏条目过滤或排序
func (n *IOSNotification) SetCategory(category string) *IOSNotification {
	n.Category = category
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

// SetOpenPage 点击打开的页面名称
func (n *WinPhoneNotification) SetOpenPage(openPage string) *WinPhoneNotification {
	n.OpenPage = openPage
	return n
}

// SetExtras 扩展字段
func (n *WinPhoneNotification) SetExtras(extras map[string]interface{}) *WinPhoneNotification {
	n.Extras = extras
	return n
}

// NewMessage 创建自定义消息实例
func NewMessage() *Message {
	return new(Message)
}

// Message 自定义消息
type Message struct {
	Content     string                 `json:"msg_content"`
	Title       string                 `json:"title,omitempty"`
	ContentType string                 `json:"content_type,omitempty"`
	Extras      map[string]interface{} `json:"extras,omitempty"`
}

// SetContent 消息内容本身
func (m *Message) SetContent(content string) *Message {
	m.Content = content
	return m
}

// SetTitle 消息标题
func (m *Message) SetTitle(title string) *Message {
	m.Title = title
	return m
}

// SetContentType 消息内容类型
func (m *Message) SetContentType(contentType string) *Message {
	m.ContentType = contentType
	return m
}

// SetExtras JSON 格式的可选参数
func (m *Message) SetExtras(extras map[string]interface{}) *Message {
	m.Extras = extras
	return m
}

// NewSmsMessage 创建短信补充实例
func NewSmsMessage() *SmsMessage {
	return new(SmsMessage)
}

// SmsMessage 短信补充
type SmsMessage struct {
	TempPara  interface{} `json:"temp_para,omitempty"`
	TempID    int64       `json:"temp_id"`
	DelayTime int         `json:"delay_time"`
}

// SetTempPara 短信模板中的参数
func (m *SmsMessage) SetTempPara(tempPara interface{}) *SmsMessage {
	m.TempPara = tempPara
	return m
}

// SetTempID 短信补充的内容模板 ID。没有填写该字段即表示不使用短信补充功能
func (m *SmsMessage) SetTempID(tempID int64) *SmsMessage {
	m.TempID = tempID
	return m
}

// SetDelayTime 单位为秒，不能超过 24 小时。设置为 0，表示立即发送短信。该参数仅对 android 和 iOS 平台有效，Winphone 平台则会立即发送短信。
func (m *SmsMessage) SetDelayTime(delayTime int) *SmsMessage {
	m.DelayTime = delayTime
	return m
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

// SetSendNO 推送序号
func (o *Options) SetSendNO(sendNO int) *Options {
	o.SendNO = sendNO
	return o
}

// SetTimeLive 离线消息保留时长(秒)
func (o *Options) SetTimeLive(timeLive int) *Options {
	o.TimeLive = timeLive
	return o
}

// SetOverrideMsgID 要覆盖的消息 ID
func (o *Options) SetOverrideMsgID(overrideMsgID int64) *Options {
	o.OverrideMsgID = overrideMsgID
	return o
}

// SetApnsCollapseID 要覆盖的消息 ID
func (o *Options) SetApnsCollapseID(apnsCollapseID string) *Options {
	o.ApnsCollapseID = apnsCollapseID
	return o
}

// SetApnsProduction 设定 APNs 是否生产环境
func (o *Options) SetApnsProduction(prod bool) *Options {
	o.ApnsProduction = prod
	return o
}

// SetBigPushDuration 定速推送时长(分钟)
func (o *Options) SetBigPushDuration(bigPushDuration int) *Options {
	o.BigPushDuration = bigPushDuration
	return o
}
