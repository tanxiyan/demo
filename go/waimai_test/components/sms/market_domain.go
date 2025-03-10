package sms

import (
	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
)

var marketSendSMS = func(siginName, note, phone, code string) {
	try.Throw(frame.CODE_WARN, "验证码发送失败:未配置服务市场验证码", "")
}

// 注册一个服务市场发送短信的配置入口
func RegMarketSendSMS(fn func(siginName, note, phone, code string)) {
	marketSendSMS = fn
}

type marketDomain struct {
	*Component
}

func (this *Component) marketDomain() *marketDomain {
	return &marketDomain{this}
}

// 创蓝短信发送
func (this *marketDomain) Send(mode, phone string, code string) {
	conf := this.Config().(*Config).Market
	marketSendSMS(conf.SignName, "安全验证", phone, code)
}
