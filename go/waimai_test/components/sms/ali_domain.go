package sms

import (
	"fmt"

	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"

	"github.com/northbright/aliyun/message"
)

// 阿里短信验证码
type aliDomain struct {
	*Component
}

func (this *Component) aliDomain() *aliDomain {
	return &aliDomain{this}
}

// 发送短信
func (this *aliDomain) Send(mode, mobile string, code string) {
	config := this.Config().(*Config)
	conf := config.Ali
	sms := message.NewClient(conf.AccessKeyID, conf.AccessKeySecret)
	code = fmt.Sprintf(`{"code":"%s"}`, code)
	ok, rsp, err := sms.SendSMS([]string{mobile}, conf.SignName, conf.TemplateCode, code)
	if !ok && rsp != nil {
		try.Throwf(frame.CODE_WARN, "验证码发送失败", "阿里，发送手机：%s，验证码：%s，错误信息：%s ", mobile, code, rsp.Message)
	}
	if err != nil {
		try.Throwf(frame.CODE_WARN, "验证码发送失败", "阿里，发送手机：%s，验证码：%s，错误信息：%s ", mobile, code, err.Error())
	}
}
