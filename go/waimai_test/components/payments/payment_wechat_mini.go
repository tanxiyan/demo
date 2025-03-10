package payments

import (
	"gitee.com/go-mao/mao/libs/utils"
	"strings"

	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
	"github.com/ArtisanCloud/PowerWeChat/v2/src/kernel/models"
	"github.com/ArtisanCloud/PowerWeChat/v2/src/payment"
	"github.com/ArtisanCloud/PowerWeChat/v2/src/payment/notify/request"
	order_request "github.com/ArtisanCloud/PowerWeChat/v2/src/payment/order/request"
	refund_request "github.com/ArtisanCloud/PowerWeChat/v2/src/payment/refund/request"
)

// 微信支付
type wechatDomain struct {
	*Component
}

func (this *Component) WechatDomain() *wechatDomain {
	object := &wechatDomain{Component: this}
	return object
}

func (this *wechatDomain) NewClient(notify string) *payment.Payment {
	wechatConfig := this.GetConfig().Wechat
	notify = this.getDomain() + notify
	client, err := payment.NewPayment(&payment.UserConfig{
		AppID:       wechatConfig.AppId,
		MchID:       wechatConfig.MchId,
		MchApiV3Key: wechatConfig.MchApiV3Key,
		Key:         wechatConfig.Key,
		CertPath:    CertPathURL,
		KeyPath:     KeyPathURL,
		SerialNo:    wechatConfig.SerialNo,
		NotifyURL:   notify,
		Http: payment.Http{
			Timeout: 30.0,
			BaseURI: "https://api.mch.weixin.qq.com",
		},
		HttpDebug: false,
		Log: payment.Log{
			Level: "debug",
			File:  "./runtime/wechat.log",
		},
	})
	if err != nil {
		try.Throw(frame.CODE_FATAL, "WechatPay初始化失败：err:", err.Error())
	}
	return client
}

func (this wechatDomain) Name() string {
	return PAYMENT_NAME_WECHAT_MINI
}

func (this wechatDomain) Type() string {
	return "wechat-jump"
}

// 是否支持充值
func (this wechatDomain) CanRecharge() bool {
	return true
}

// 是否支持提现
func (this wechatDomain) CanWithdraw() bool {
	return true
}

func (this wechatDomain) CanAutoWithdraw() bool {
	return true
}

func (this wechatDomain) UseMoneyValid() bool {
	return false
}

func (this *wechatDomain) getDomain() string {
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

// PayWeb
//
//	@Description: 网页支付
//	@receiver this
//	@param body
//	@return jumpURL
func (this *wechatDomain) PayWeb(body PayBody) (res interface{}) {
	client := this.NewClient(body.NotifyURL)
	options := &order_request.RequestJSAPIPrepay{
		Amount: &order_request.JSAPIAmount{
			Total:    int(body.Amount * 100),
			Currency: "CNY",
		},
		Attach:      body.TradeNO,
		Description: body.Subject,
		OutTradeNo:  body.TradeNO,
		Payer: &order_request.JSAPIPayer{
			OpenID: body.OpenId,
		},
	}
	native, err := client.Order.JSAPITransaction(options)
	if err != nil {
		try.Throw(frame.CODE_FATAL, "下单失败，err：", err.Error())
		return ""
	}
	payConf, err := client.JSSDK.BridgeConfig(native.PrepayID, true)
	if err != nil {
		try.Throw(frame.CODE_FATAL, "调起支付失败，err：", err.Error())
		return ""
	}
	return payConf
}

// Transfer
//
//	@Description: 商户转账
//	@receiver this
//	@param body
//	@return outTradeNO
func (this *wechatDomain) Transfer(body TransferBody) (outTradeNO string) {
	client := this.NewClient("")
	options := &refund_request.RequestRefund{
		OutTradeNo:   body.TradeNo,
		OutRefundNo:  utils.MakeSN(),
		Reason:       body.ShowName,
		FundsAccount: "",
		Amount: &refund_request.RefundAmount{
			Refund:   int(body.Amount * 100),               // 退款金额，单位：分
			Total:    int(body.Amount * 100),               // 订单总金额，单位：分
			From:     []*refund_request.RefundAmountFrom{}, // 退款出资账户及金额。不传仍然需要这个空数组防止微信报错
			Currency: "CNY",
		},
		GoodsDetail: nil,
	}
	response, err := client.Refund.Refund(options)
	if err != nil {
		try.Throw(frame.CODE_FATAL, "转账失败，err：", err.Error())
		return ""
	}
	if response.Message != "" {
		try.Throw(frame.CODE_FATAL, "转账失败，err：", response.Message)
		return ""
	}
	return response.OutTradeNO
}

// CheckNotification
//
//	@Description: 回调验证
//	@receiver this
//	@param sn
//	@param c
//	@return res
//	@return outerTradeNo
func (this *wechatDomain) CheckNotification(sn string, w *frame.Webline) (res, outerTradeNo string) {
	client := this.NewClient("")
	response, err := client.HandlePaidNotify(w.Request,
		func(message *request.RequestNotify, transaction *models.Transaction, fail func(message string)) interface{} {
			// 看下支付通知事件状态
			// 这里可能是微信支付失败的通知，所以可能需要在数据库做一些记录，然后告诉微信我处理完成了。
			if message.EventType != "TRANSACTION.SUCCESS" {
				return true
			}
			if transaction.OutTradeNo == "" {
				// 因为微信这个回调不存在订单号，所以可以告诉微信我还没处理成功，等会它会重新发起通知
				// 如果不需要，直接返回true即可
				fail("payment fail")
				return nil
			}
			outerTradeNo = transaction.OutTradeNo
			return true
		},
	)
	// 这里可能是因为不是微信官方调用的，无法正常解析出transaction和message，所以直接抛错。
	if err != nil {
		try.Throw(frame.CODE_FATAL, "微信回调解析失败，err：", err.Error())
	}
	// 这里根据之前返回的是true或者fail，框架这边自动会帮你回复微信
	err = response.Write(w.Writer)
	if err != nil {
		try.Throw(frame.CODE_FATAL, "微信回调返回失败，err：", err.Error())
	}
	return "success", outerTradeNo
}
