package payments

import (
	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
)

type Component struct {
	config *Config
	*frame.Taskline
	notifyDomain string
}

func New(t *frame.Taskline) *Component {
	object := new(Component)
	object.Taskline = t
	return object
}

func (this *Component) CompName() string {
	return "Payments"
}

func (this *Component) Config() frame.ConfigInterface {
	return this.RenderConfig(&Config{})
}

func (this *Component) GetConfig() *Config {
	return this.Config().(*Config)
}

func (this *Component) Models() []frame.ModelInterface {
	return []frame.ModelInterface{}
}

func (this *Component) NewPayment(paymentType string) PaymentInterface {
	switch paymentType {
	case PAYMENT_NAME_WECHAT_MINI:
		return this.WechatDomain()
	default:
		try.Throwf(frame.CODE_WARN, "不支持%s支付", paymentType)
	}
	return nil
}
