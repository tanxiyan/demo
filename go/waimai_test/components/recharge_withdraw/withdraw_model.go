package recharge_withdraw

import (
	"gitee.com/go-mao/mao/frame"
	"time"
)

const (
	FETCH_STATE_SQ = 1  //处理中
	FETCH_STATE_CG = 2  //成功
	FETCH_STATE_SB = -1 //失败
	FETCH_STATE_TH = -2 //退回
)

// 提现
type WithdrawModel struct {
	Id                 int64
	UserId             int64         `xorm:"int(11) index comment('用户id')"`
	SN                 string        `xorm:"varchar(16) unique comment('sn')"`
	State              int           `xorm:"tinyint(4) comment('提现状态'')"`
	ServiceCharge      int64         `xorm:"bigint(20) comment('手续费')"`
	Total              int64         `xorm:"bigint(20) index comment('扣除总量=提取数量+手续费')"`
	Number             int64         `xorm:"bigint(20) comment('金额，单位：元')"`
	Remark             string        `xorm:"varchar(255) comment('备注')"`
	PaymentAccount     string        `xorm:"varchar(64) comment('提现账户')"`
	PaymentAccountType string        `xorm:"varchar(24) comment('提现账户类型')"`
	PaymentQrcode      string        `xorm:"varchar(255) comment('收款码')"`
	RealName           string        `xorm:"varchar(32) comment('提现真实姓名')"`
	IP                 string        `xorm:"varchar(32) comment('提现ip')"`
	OuterOrderId       string        `xorm:"varchar(32) comment('外部交易订单')"`
	WithdrawLogs       []WithdrawLog `xorm:"json comment('日志记录')"`
	UpdateBy           int64         `xorm:"bigint(20)"`
	CreateTime         time.Time     `xorm:"created datetime comment('创建时间')"`
	UpdateTime         time.Time     `xorm:"updated datetime comment('创建时间')"`
	prefix             string        `xorm:"-"`
	frame.Model        `xorm:"-"`
}

// 错误记录
type WithdrawLog struct {
	CreateTime string `json:"CreateTime"` //日志时间
	Info       string `json:"info"`       //说明
}

func (this *Component) NewWithdraw() *WithdrawModel {
	withdraw := this.OrmModel(&WithdrawModel{prefix: this.prefix}).(*WithdrawModel)
	withdraw.SetTableName(this.prefix + "_withdraw")
	return withdraw
}

// 设置主键值（
func (this *WithdrawModel) PrimaryKey() interface{} {
	return this.Id
}
