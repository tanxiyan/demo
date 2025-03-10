package recharge_withdraw

import (
	"errors"
	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
	"gitee.com/go-mao/mao/libs/utils"
	"time"
	"waimai_api/components/account"
	"waimai_api/components/coin"
	"waimai_api/components/payment_account"
)

type withdrawDomain struct {
	*Component
}

func (this *Component) WithdrawDomain() *withdrawDomain {
	return &withdrawDomain{this}
}

// 发起提现
func (this *withdrawDomain) Withdraw(user account.UserInterface, paymentAccount *payment_account.AccountModel, number coin.Coin, createIP string) (r *WithdrawModel) {
	if int64(number)%coin.GetRateByCoinType(coin.GetCoinType()-2) != 0 {
		try.Throw(frame.CODE_WARN, "提现金额小数位数不能超过2位")
		return
	}
	conf := this.GetConfig()
	todayWithdrawTimes := this.CountWithdrawTimes(user.GetUserId(), time.Now().Format(utils.FORMAT_DATE_TIME), FETCH_STATE_SQ, FETCH_STATE_CG)
	if todayWithdrawTimes >= conf.WithdrawOneDay {
		try.Throwf(frame.CODE_WARN, "每日提现至多%d次", conf.WithdrawOneDay)
		return
	}
	serviceNumber := coin.CoinByCharge(number, conf.WithdrawCharge)
	//
	withdraw := this.NewWithdraw()
	withdraw.UserId = user.GetUserId()
	withdraw.SN = utils.MakeSN()
	withdraw.ServiceCharge = serviceNumber
	withdraw.Number = int64(number)
	withdraw.Total = withdraw.ServiceCharge + withdraw.Number
	withdraw.State = FETCH_STATE_SQ
	withdraw.PaymentAccountType = paymentAccount.AccountType
	withdraw.RealName = paymentAccount.Realname
	withdraw.PaymentAccount = paymentAccount.Account
	withdraw.PaymentQrcode = paymentAccount.Qrcode
	withdraw.IP = createIP
	//提现总数
	if withdraw.Number < int64(conf.LowWithdrawMoney) {
		try.Throwf(frame.CODE_WARN, "单笔最低提现金额为%s元", conf.LowWithdrawMoney)
		return
	}
	if withdraw.Number > int64(conf.MaxWithdrawMoney) {
		try.Throwf(frame.CODE_WARN, "单笔最大提现金额为%s元", conf.MaxWithdrawMoney)
		return
	}
	this.Transaction(func() {
		if !withdraw.Create() {
			try.Throwf(frame.CODE_SQL, "提现申请失败%d", user.GetUserId())
		}
		this.handler.WithdrawBefore(this.Taskline, withdraw)
		if withdraw.Total > int64(conf.AutoWithdrawMoney) {
			return
		}
		if !this.handler.CanAutoWithdraw(this.Taskline, withdraw) {
			return
		}
		//调用支付接口
		_ = this.withdrawHandle(withdraw)
	})
	return withdraw
}

// 审核退回
func (this *withdrawDomain) CheckBack(user UserInterface, withdrawSN string, remark string) {
	if remark == "" {
		try.Throw(frame.CODE_WARN, "请输入退回原因")
		return
	}
	withdraw, ok := this.GetBySN(withdrawSN)
	if !ok {
		try.Throw(frame.CODE_WARN, "无效的提现订单")
		return
	}
	if !utils.ArrayIn(withdraw.State, FETCH_STATE_SQ, FETCH_STATE_SB) {
		try.Throw(frame.CODE_WARN, "当前状态无法操作")
		return
	}
	this.Transaction(func() {
		withdraw.State = FETCH_STATE_TH
		withdraw.Remark = remark
		withdraw.UpdateBy = user.GetUserId()
		this.handler.WithdrawBack(this.Taskline, withdraw)
		if !withdraw.Cols("State", "Remark").Update() {
			try.Throw(frame.CODE_SQL, "操作失败")
		}
	})
}

// 审核通过
func (this *withdrawDomain) CheckPass(user UserInterface, withdrawSN string, createIP string) {
	withdraw, ok := this.GetBySN(withdrawSN)
	if !ok {
		try.Throw(frame.CODE_WARN, "无效的数据")
		return
	}
	if !utils.ArrayIn(withdraw.State, FETCH_STATE_SQ, FETCH_STATE_SB) {
		try.Throw(frame.CODE_WARN, "当前状态无法操作")
		return
	}
	withdraw.UpdateBy = user.GetUserId()
	var err error
	this.Transaction(func() {
		withdraw.IP = createIP
		err = this.withdrawHandle(withdraw)
	})
	if err != nil {
		try.Throw(frame.CODE_FATAL, err.Error())
	}
}

// 调用支付接口发起提现(外部提供事务）
func (this *withdrawDomain) withdrawHandle(withdraw *WithdrawModel) (err error) {
	var outerOrderId string
	try.Do(func() {
		outerOrderId = this.handler.WithdrawAPI(this.Taskline, withdraw)
		withdraw.OuterOrderId = outerOrderId
		withdraw.State = FETCH_STATE_CG
		withdraw.Cols("OuterOrderId", "State")
	}, func(e try.Exception) {
		withdraw.WithdrawLogs = append(withdraw.WithdrawLogs, WithdrawLog{
			CreateTime: time.Now().Format(utils.FORMAT_DATE_TIME),
			Info:       e.ErrMsg(),
		})
		withdraw.Cols("WithdrawLogs")
		err = errors.New(e.ErrMsg())
	})
	if !withdraw.Update() {
		try.Throw(frame.CODE_SQL, "提现接口调用后数据保存失败")
		return
	}
	if withdraw.State == FETCH_STATE_CG {
		this.handler.WithdrawSucc(this.Taskline, withdraw)
	}
	return
}

type FindWithdrawArgs struct {
	UserId  int64
	BeginAt string
	EndAt   string
	State   int
	SN      string
}

// 提现记录
func (this *withdrawDomain) Find(args FindWithdrawArgs, page, limit int, listPtr interface{}) int64 {
	db := this.OrmTable(this.NewWithdraw())
	if args.UserId > 0 {
		db.And("`UserId`=?", args.UserId)
	}
	if args.BeginAt != "" {
		db.And("`CreateTime`>?", args.BeginAt)
	}
	if args.EndAt != "" {
		db.And("`CreateTime`<?", args.EndAt)
	}
	if args.State != 0 {
		db.And("`State`=?", args.State)
	}
	if args.SN != "" {
		db.And("`SN` like ?", "%"+args.SN+"%")
	}
	db.Desc("Id")
	return this.FindPage(db, listPtr, page, limit)
}

func (this *withdrawDomain) GetById(id int64) (*WithdrawModel, bool) {
	withdraw := this.NewWithdraw()
	ok := withdraw.Match("Id", id).Get()
	return withdraw, ok
}

func (this *withdrawDomain) GetBySN(sn string) (*WithdrawModel, bool) {
	withdraw := this.NewWithdraw()
	ok := withdraw.Match("SN", sn).Get()
	return withdraw, ok
}

// 提现次数统计
func (this *withdrawDomain) CountWithdrawTimes(userId int64, beginAt string, states ...interface{}) int64 {
	db := this.OrmTable(this.NewWithdraw()).Where("`UserId`=?", userId)
	if len(states) > 0 {
		db.In("State", states...)
	}
	if beginAt != "" {
		db.And("`CreateTime`>?", beginAt)
	}
	total, err := db.Count()
	if err != nil {
		try.Throw(frame.CODE_SQL, "提现次数统计失败，err：", err.Error())
	}
	return total
}

// 提现总额(只包含提现成功的）
func (this *withdrawDomain) Sum(userId int64, beginDate, endDate string) int64 {
	withdrawInfo := this.NewWithdraw()
	db := this.OrmTable(withdrawInfo).Where("`State`=?", FETCH_STATE_CG)
	if userId > 0 {
		db.And("`UserId`=?", userId)
	}
	if beginDate != "" {
		db.And("`CreateTime`>?", beginDate)
	}
	if endDate != "" {
		db.And("`CreateTime`<?", endDate)
	}
	total, err := db.SumInt(withdrawInfo, "Number")
	if err != nil {
		try.Throw(frame.CODE_SQL, "提现总额统计失败，err：", err.Error())
	}
	return total
}

func (this *withdrawDomain) DeleteByIds(userIds []int64) {
	_, err := this.OrmTable(this.NewWithdraw()).In("UserId", userIds).Delete()
	if err != nil {
		try.Throw(frame.CODE_SQL, "提现记录清理失败", err.Error())
	}
}
