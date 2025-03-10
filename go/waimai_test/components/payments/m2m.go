package payments

import (
	"fmt"
	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay"
	"strings"
	"waimai_api/globals"
)

type m2mDomain struct {
	*Component
	alipayClient *alipay.AliPay
}

func (this *Component) M2MDomain() *m2mDomain {
	object := &m2mDomain{Component: this}
	config := this.GetConfig().M2MConfig
	object.alipayClient = alipay.New(config.AppId, config.PublicKey, config.PrivateKey, true)
	return object
}

func (this m2mDomain) Name() string {
	return PAYMENT_NAME_M2M
}

func (this m2mDomain) Type() string {
	return "alipay-qrcode"
}

// 是否支持充值
func (this m2mDomain) CanRecharge() bool {
	return false
}

// 是否支持提现
func (this m2mDomain) CanWithdraw() bool {
	return false
}

func (this m2mDomain) CanAutoWithdraw() bool {
	return false
}

func (this m2mDomain) UseMoneyValid() bool {
	return false
}

func (this *m2mDomain) getDomain() string {
	config := this.GetConfig()
	domain := this.notifyDomain
	if domain == "" {
		domain = config.ApiDomainName
	}
	if strings.HasPrefix(domain, "http") {
		return domain
	}
	return "http://" + domain
}

// 网页支付
func (this *m2mDomain) PayWeb(body PayBody) (jumpURL string) {
	aliClient := this.alipayClient
	p := alipay.AliPayTradePreCreate{}
	p.OutTradeNo = body.TradeNO
	p.TotalAmount = fmt.Sprintf("%.2f", body.Amount)
	p.Subject = body.Subject
	p.StoreId = "NO001"
	p.NotifyURL = this.getDomain() + body.NotifyURL
	p.TimeoutExpress = "60m"
	r, err := aliClient.TradePreCreate(p)
	if err != nil {
		try.Throwf(frame.CODE_FATAL, "发起充值失败%s", err.Error())
	}
	if !r.IsSuccess() || r.AliPayPreCreateResponse.QRCode == "" {
		try.Throwf(frame.CODE_FATAL, "发起充值失败，%s%s", r.AliPayPreCreateResponse.Msg, r.AliPayPreCreateResponse.SubMsg)
	}
	return this.getDomain() + globals.MakeQrcodeLink(r.AliPayPreCreateResponse.QRCode)
}

func (this *m2mDomain) Transfer(body TransferBody) (outTradeNO string) {
	try.Throw(frame.CODE_WARN, "不支持的支付方式")
	return ""
}

// 交易通知验证
func (this *m2mDomain) CheckNotification(sn string, c *gin.Context) (res, outerTradeNo string) {
	notify, err := this.alipayClient.GetTradeNotification(c.Request)
	if err != nil {
		try.Throwf(frame.CODE_FATAL, "支付宝充值通知内容错误，单号：%s 错误：%s", sn, err.Error())
	}
	if notify.TradeStatus != "TRADE_SUCCESS" { //交易失败
		try.Throwf(frame.CODE_FATAL, "支付宝充值交易失败，单号：%s", sn)
	}
	return "success", notify.TradeNo
}
