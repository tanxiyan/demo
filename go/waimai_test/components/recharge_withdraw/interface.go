package recharge_withdraw

import (
	"waimai_api/components/payments"

	"gitee.com/go-mao/mao/frame"
)

type RechargeWithdrawhHandlerInteface interface {
	Prefix() string
	GetConfig() *Config
	WithdrawAPI(b *frame.Taskline, withdraw *WithdrawModel) string                                                        //提现处理器（根据支付方式获取提现方式）
	WithdrawBefore(b *frame.Taskline, withdraw *WithdrawModel)                                                            //提现前(可扣款）
	WithdrawSucc(b *frame.Taskline, withdraw *WithdrawModel)                                                              //提现成功后（可发起通知）
	WithdrawBack(b *frame.Taskline, withdraw *WithdrawModel)                                                              //提现退回（可退款）
	CanAutoWithdraw(b *frame.Taskline, withdraw *WithdrawModel) bool                                                      //是否可以自动提现
	UseMoneyValid(b *frame.Taskline, rechargeType string) bool                                                            //是否需要金额自动增加处理
	RechargeAPI(b *frame.Taskline, withdraw *RechargeModel, webPayMode payments.WebPayMode, returnURL string) interface{} //充值处理,返回支付url
	CheckNotification(b *frame.Taskline, recharge *RechargeModel, w *frame.Webline) (notify_res, outerTradeNo string)
	RechargeSucc(b *frame.Taskline, recharge *RechargeModel) //充值成功（可加款）
}

// 用户接口
type UserInterface interface {
	GetUserId() int64
	GetUserGroup() string
}
