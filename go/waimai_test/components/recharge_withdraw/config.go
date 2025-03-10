package recharge_withdraw

import (
	"waimai_api/components/coin"
)

type Config struct {
	MinRecharge       coin.Coin `name:"最低充值金额（元）"`
	WithdrawShowName  string    /*`name:"提现标题" desc:"例如：xx酒店" valid:"required->提现标题不能为空"`*/
	RechargeSubject   string    `name:"充值标题" desc:"例如：订购客房" valid:"required->充值标题不能为空"`
	WithdrawOneDay    int64     /*`name:"每日可提现次数" valid:"gte=0->每日可提现次数不能小于0"`*/
	WithdrawCharge    float64   /*`name:"提现手续费率/%" valid:"gte=0,lt=100->提现手续费率在0-100之间"`*/
	LowWithdrawMoney  coin.Coin /*`name:"最低提现金额" valid:"gt=0->最低提现金额不能小于0" desc:"单位：元"`*/
	AutoWithdrawMoney coin.Coin /*`name:"自动提现最大金额" valid:"gt=0->自动提现最大金额不能小于0" desc:"低于该金额自动提现"`*/
	MaxWithdrawMoney  coin.Coin /*`name:"单笔最大可提现金额" valid:"gt=0->单笔最大提现金额不能小于0" desc:"单位：元"`*/
}
