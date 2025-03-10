package sms

// 发送短信接口
type Sender interface {
	Send(mode string, mobile, code string)
}
