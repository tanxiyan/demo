package recharge_withdraw

import (
	"gitee.com/go-mao/mao/frame"
	"time"
)

const (
	RECHARGE_STATE_DFK = 1  //待付款
	RECHARGE_STATE_YFK = 2  //已付款
	RECHARGE_STATE_YGQ = -2 //已过期
	RECHARGE_STATE_SB  = -1 //失败
	RECHARGE_STATE_TH  = -3 // 已退回
)

// 充值
type RechargeModel struct {
	Id           int64
	SN           string        `xorm:"varchar(16) unique comment('单号')"`
	ShortSN      string        `xorm:"varchar(8) unique comment('短码单号')"`
	UserId       int64         `xorm:"int(11) index comment('用户id')"`
	Number       int64         `xorm:"bigint(20) index comment('数量')"`
	State        int           `xorm:"tinyint(4) comment('付款状态')"`
	WebPayMode   int           `xorm:"tinyint(4) comment('支付类型1移动端、2pc端')"`
	OuterTradeSN string        `xorm:"varchar(64) index comment('外部支付单号')"`
	RechargeType string        `xorm:"varchar(40) comment('充值方式')"`
	BaseNumber   int64         `xorm:"bigint(20) comment('基础充值数量')"`
	Relation     string        `xorm:"varchar(255) comment('关联数据')"`
	Remarks      string        `xorm:"text comment('备注')"`
	RechargeLogs []RechargeLog `xorm:"json comment('日志记录')"`
	Version      int           `xorm:"version"`
	UpdateBy     int64         `xorm:"bigint(20)"`
	CreateTime   time.Time     `xorm:"created datetime comment('创建时间')"`
	UpdateTime   time.Time     `xorm:"updated datetime comment('创建时间')"`
	ArrivalTime  time.Time     `xorm:"datetime comment('到账时间')"`
	prefix       string        `xorm:"-"`
	frame.Model  `xorm:"-"`
}

type RechargeLog struct {
	CreateTime string `json:"CreateTime"` //日志时间
	Info       string `json:"info"`       //说明
}

func (this *Component) NewRecharge() *RechargeModel {
	recharge := this.OrmModel(&RechargeModel{prefix: this.prefix}).(*RechargeModel)
	recharge.SetTableName(this.prefix + "_recharge")
	return recharge
}

// 设置主键值（
func (this *RechargeModel) PrimaryKey() interface{} {
	return this.Id
}
