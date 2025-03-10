package recharge_withdraw

import (
	"sync"
	"time"
	"waimai_api/components/account"
	"waimai_api/components/coin"
	"waimai_api/components/payments"

	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
	"gitee.com/go-mao/mao/libs/utils"
)

type rechargeDomain struct {
	*Component
}

func (this *Component) RechargeDomain() *rechargeDomain {
	return &rechargeDomain{this}
}

// 根据有效金额数量查询,如果当前时间线没有，则查询上一次的
func (this *rechargeDomain) GetByValidNumber(number int64, rechargeType string) (*RechargeModel, bool) {
	rechargeInfo := this.NewRecharge()
	fiveMinutesAgo := time.Now().Add(0 - time.Second*300).Format(utils.FORMAT_DATE_TIME)
	rechargeInfo.Where("`Number`=? and `State`=? and `CreateTime`>? and `RechargeType`=?", number, RECHARGE_STATE_DFK, fiveMinutesAgo, rechargeType)
	ok := rechargeInfo.Get()
	return rechargeInfo, ok
}

func (this *rechargeDomain) GetByRelation(relation string) (*RechargeModel, bool) {
	rechargeInfo := this.NewRecharge()
	has := rechargeInfo.Match("Relation", relation).Get()
	return rechargeInfo, has
}

//mode=充值方式,如果5分钟内有相同的未付款充值单则不再创建新充值单
/*
时间限定=5分钟
1、如果A充值1元，当前没有其他人充1元，则只需要充1.00元
2、如果B需要充1元，则需要充1.01元，自动生成
3、如果C直接输入充值1.01元，当前有人占用了，则累加，需要充1.02元，自动生成
4、如果是游客下单订单总额为1.00元，这个游客需要付款1.03元才能正常下单
5、叠加数量不能超过100，也就是同一金额5分钟内最多可跨度100分，超过100分还有错误的会报“当前付款人数过多，请稍后再试”
*/
//relation=关联参数，可以为空
var rechargeLock sync.Mutex

func (this *rechargeDomain) Recharge(user account.UserInterface, baseNumber int64, rechargeType string, webPayMode payments.WebPayMode, returnURL, relation, remark string) (rechargeInfo *RechargeModel, payURL interface{}) {
	rechargeLock.Lock()
	defer rechargeLock.Unlock()
	if baseNumber <= 0 {
		try.Throw(frame.CODE_WARN, "金额不能小于0")
		return
	}
	coinZeroNum := coin.GetRateByCoinType(coin.GetCoinType() - 2) //
	if baseNumber%coinZeroNum != 0 {
		try.Throw(frame.CODE_WARN, "金额小数位数不能超过2位")
		return
	}
	rechargeInfo = this.NewRecharge()
	rechargeInfo.Number = baseNumber
	this.Transaction(func() {
		rechargeInfo.Relation = relation
		rechargeInfo.UserId = user.GetUserId()
		rechargeInfo.BaseNumber = baseNumber
		rechargeInfo.SN = utils.MakeSN()
		rechargeInfo.State = RECHARGE_STATE_DFK
		rechargeInfo.RechargeType = rechargeType
		logs := make([]RechargeLog, 0)
		logs = append(logs, RechargeLog{
			CreateTime: time.Now().Format(utils.FORMAT_DATE_TIME),
			Info:       remark,
		})
		rechargeInfo.RechargeLogs = logs
		rechargeInfo.Remarks = remark
		if !rechargeInfo.Create() {
			try.Throwf(frame.CODE_SQL, "用户%d充值日志保存失败", user.GetUserId())
			return
		}
		rechargeInfo.ShortSN = utils.MakeShortSN(rechargeInfo.Id + 10000)
		if !rechargeInfo.Cols("ShortSN").Update() {
			try.Throw(frame.CODE_SQL, "充值失败，短码修改失败")
			return
		}
		payURL = this.handler.RechargeAPI(this.Taskline, rechargeInfo, webPayMode, returnURL)
	})
	return
}

// 充值通知处理,res=返回给通知方的数据
func (this *rechargeDomain) RechargeNotify(handlerUserId int64, sn string, w *frame.Webline) (res string) {
	rechargeInfo, ok := this.GetBySN(sn)
	if !ok {
		try.Throw(frame.CODE_WARN, "交易失败，订单无效")
		return
	}
	var outerTradeSN string
	if handlerUserId > 0 {
		outerTradeSN = ""
	} else {
		res, outerTradeSN = this.handler.CheckNotification(this.Taskline, rechargeInfo, w)
		if rechargeInfo.State == RECHARGE_STATE_YFK {
			return res
		}
		hasOuterTradeTotal, _ := this.OrmTable(this.NewRecharge()).Where("`OuterTradeSN`=? and `RechargeType`=?", outerTradeSN, rechargeInfo.RechargeType).Count()
		if hasOuterTradeTotal > 0 {
			return
		}
	}
	if rechargeInfo.State != RECHARGE_STATE_DFK {
		return
	}
	this.Transaction(func() {
		rechargeInfo.UpdateBy = handlerUserId
		rechargeInfo.OuterTradeSN = outerTradeSN
		rechargeInfo.State = RECHARGE_STATE_YFK
		rechargeInfo.ArrivalTime = time.Now()
		if !rechargeInfo.Cols("State", "ArrivalTime", "OuterTradeSN", "UpdateBy").Update() {
			try.Throw(frame.CODE_SQL, "充值失败")
		}
		this.handler.RechargeSucc(this.Taskline, rechargeInfo)
	})
	return res
}

type RechargeFindArgs struct {
	UserId  int64
	State   int
	BeginAt string
	EndAt   string
}

// 出售记录查询
func (this *rechargeDomain) Find(args *RechargeFindArgs, page, limit int, listPtr interface{}) int64 {
	db := this.OrmTable(this.NewRecharge())
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
	db.Desc("Id")
	return this.FindPage(db, listPtr, page, limit)
}

func (this *rechargeDomain) GetById(id int64) (*RechargeModel, bool) {
	rechargeInfo := this.NewRecharge()
	ok := rechargeInfo.Match("Id", id).Get()
	return rechargeInfo, ok
}

func (this *rechargeDomain) GetBySN(sn string) (*RechargeModel, bool) {
	rechargeInfo := this.NewRecharge()
	ok := rechargeInfo.Match("SN", sn).Get()
	return rechargeInfo, ok
}

func (this *rechargeDomain) GetByShortSN(sn string) (*RechargeModel, bool) {
	rechargeInfo := this.NewRecharge()
	ok := rechargeInfo.Match("ShortSN", sn).Get()
	return rechargeInfo, ok
}

// 充值总量
func (this *rechargeDomain) Sum(beginDate, endDate string) int64 {
	rechargeInfo := this.NewRecharge()
	db := this.OrmTable(rechargeInfo).Where("`State`=?", RECHARGE_STATE_YFK)
	if beginDate != "" {
		db.And("`CreateTime`>?", beginDate)
	}
	if endDate != "" {
		db.And("`CreateTime`<?", endDate)
	}
	total, err := db.SumInt(rechargeInfo, "Number")
	if err != nil {
		try.Throw(frame.CODE_SQL, "统计失败")
	}
	return total
}

func (this *rechargeDomain) DeleteByIds(userIds []int64) {
	_, err := this.OrmTable(this.NewRecharge()).In("UserId", userIds).Delete()
	if err != nil {
		try.Throw(frame.CODE_SQL, "充值记录清理失败", err.Error())
	}
}
