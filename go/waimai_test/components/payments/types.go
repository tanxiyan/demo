package payments

import (
	"gitee.com/go-mao/mao/frame"
)

type WebPayMode int

const (
	WEBPAY_MODE_WAP WebPayMode = 1 //移动端
	WEBPAY_MODE_PC  WebPayMode = 2 //pc端
)

// 提现参数
type TransferBody struct {
	Account  string
	Realname string
	TradeNo  string
	Amount   float64
	Remark   string
	CreateIP string
	ShowName string
}

// 支付参数
type PayBody struct {
	WebPayMode WebPayMode
	UserId     int64
	Subject    string
	TradeNO    string
	ShortSN    string
	Amount     float64
	NotifyURL  string
	ReturnURL  string
	OpenId     string
}

// 支付接口
type PaymentInterface interface {
	Name() string                                                                    //支付方式名称
	Type() string                                                                    //支付账号类型
	CanRecharge() bool                                                               //是否支持充值
	CanWithdraw() bool                                                               //是否支持提现
	CanAutoWithdraw() bool                                                           //是否支持自动提现
	UseMoneyValid() bool                                                             //使用金额验证方式
	PayWeb(body PayBody) (res interface{})                                           //付款
	Transfer(body TransferBody) (outTradeNO string)                                  //收款
	CheckNotification(sn string, w *frame.Webline) (res string, outerTradeSN string) //交易通知验证（一般可不用）
}

const (
	PAYMENT_NAME_ALIPAY      = "支付宝"
	PAYMENT_NAME_WECHAT_MINI = "微信小程序"
	PAYMENT_NAME_M2M         = "支付宝扫码"
)

const (
	PAYMENT_TYPE_ACCOUNT = "account"
	PAYMENT_TYPE_QRCODE  = "qrcode"
)
