package payments

import (
	"fmt"
	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay"
)

// 阿里支付
type alipayDomain struct {
	*Component
	alipayClient *alipay.AliPay
}

func (this *Component) AlipayDomain() *alipayDomain {
	object := &alipayDomain{Component: this}
	config := this.GetConfig().Alipay
	object.alipayClient = alipay.New(config.AppId, config.PublicKey, config.PrivateKey, true)
	return object
}

func (this *Component) AlipayDomainByConfig(config *AlipayConfig) *alipayDomain {
	object := &alipayDomain{Component: this}
	object.alipayClient = alipay.New(config.AppId, config.PublicKey, config.PrivateKey, true)
	return object
}

func (this alipayDomain) Name() string {
	return PAYMENT_NAME_ALIPAY
}

func (this alipayDomain) Type() string {
	return "alipay-jump"
}

// 是否支持充值
func (this alipayDomain) CanRecharge() bool {
	return true
}

// 是否支持提现
func (this alipayDomain) CanWithdraw() bool {
	return true
}

func (this alipayDomain) CanAutoWithdraw() bool {
	return true
}

func (this alipayDomain) UseMoneyValid() bool {
	return false
}

func (this *alipayDomain) getDomain() string {
	domain := this.notifyDomain
	if strings.HasPrefix(domain, "http") {
		return domain
	}
	return "http://" + domain
}

// 网页支付
func (this *alipayDomain) PayWeb(body PayBody) (jumpURL string) {
	var u *url.URL
	var e error
	config := this.GetConfig()
	if config.ApiDomainName != "" {
		this.notifyDomain = config.ApiDomainName
	}
	if body.WebPayMode == WEBPAY_MODE_WAP {
		p := alipay.AliPayTradeWapPay{}
		p.NotifyURL = this.getDomain() + body.NotifyURL
		p.ReturnURL = body.ReturnURL
		p.Subject = body.Subject
		p.OutTradeNo = body.TradeNO
		p.TotalAmount = fmt.Sprintf("%.2f", body.Amount)
		p.GoodsType = "0"
		p.ProductCode = "QUICK_WAP_WAY"
		u, e = this.alipayClient.TradeWapPay(p)
	} else {
		p := alipay.AliPayTradePagePay{}
		p.NotifyURL = this.getDomain() + body.NotifyURL
		p.ReturnURL = body.ReturnURL
		p.Subject = body.Subject
		p.OutTradeNo = body.TradeNO
		p.TotalAmount = fmt.Sprintf("%.2f", body.Amount)
		p.GoodsType = "0"
		p.ProductCode = "FAST_INSTANT_TRADE_PAY"
		u, e = this.alipayClient.TradePagePay(p)
	}
	if e != nil {
		try.Throwf(frame.CODE_FATAL, "支付宝充值失败，错误：%s", e.Error())
	}
	return u.String()
}

// 转账
func (this *alipayDomain) Transfer(body TransferBody) (outTradeNO string) {
	p := alipay.AliPayFundTransToAccountTransfer{}
	p.Amount = fmt.Sprintf("%.2f", body.Amount)
	p.OutBizNo = body.TradeNo       // string `json:"out_biz_no"`      // 必选 商户转账唯一订单号
	p.PayeeType = "ALIPAY_LOGONID"  // string `json:"payee_type"`      // 必选 收款方账户类型,"ALIPAY_LOGONID":支付宝帐号
	p.PayeeAccount = body.Account   // string `json:"payee_account"`   // 必选 收款方账户。与payee_type配合使用
	p.PayerShowName = body.ShowName // string `json:"payer_show_name"` // 可选 付款方显示姓名
	p.PayeeRealName = body.Realname // string `json:"payee_real_name"` // 可选 收款方真实姓名,如果本参数不为空，则会校验该账户在支付宝登记的实名是否与收款方真实姓名一致。
	p.Remark = body.Remark          // string `json:"remark"`          // 可选 转账备注,金额大于50000时必填
	rsp, err := this.alipayClient.FundTransToAccountTransfer(p)
	if err != nil {
		try.Throwf(frame.CODE_FATAL, "支付宝提现失败，err：%s", err.Error())
	}
	if !rsp.IsSuccess() {
		try.Throwf(frame.CODE_FATAL, "支付宝提现失败，错误代码：%s,错误原因：%s", rsp.Body.Code, rsp.Body.SubMsg)
	}
	return rsp.Body.OrderId
}

// 交易通知验证
func (this *alipayDomain) CheckNotification(sn string, c *gin.Context) (res, outerTradeNo string) {
	notify, err := this.alipayClient.GetTradeNotification(c.Request)
	if err != nil {
		try.Throwf(frame.CODE_FATAL, "支付宝充值通知内容错误，单号：%s 错误：%s", sn, err.Error())
	}
	if notify.TradeStatus != "TRADE_SUCCESS" { //交易失败
		try.Throwf(frame.CODE_FATAL, "支付宝充值交易失败，单号：%s", sn, err.Error())
	}
	return "success", notify.TradeNo
}
