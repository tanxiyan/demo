package account

import (
	"time"

	"gitee.com/go-mao/mao/frame"
)

type AccountLogModel struct {
	Id              int64
	UserId          int64     `xorm:"notnull int(11) index comment('用户id')"`
	Number          int64     `xorm:"bigint(20) comment('数量')"`
	BalanceBefore   int64     `xorm:"bigint(20) comment('之前数量')"`
	TypeId          int       `xorm:"tinyint(4) comment('日志类型id')"`
	Remarks         string    `xorm:"varchar(255) comment('日志备注【管理员可见】')"`
	Note            string    `xorm:"varchar(255) comment('备注【会员可见】')"`
	BalanceAfter    int64     `xorm:"bigint(20) comment('之后数量')"`
	CreateUserId    int64     `xorm:"bigint(20) comment('创建用户')"`
	CreateUserGroup string    `xorm:"varchar(32) comment('创建用户所在组')"`
	CreateTime      time.Time `xorm:"created comment('创建时间')"`
	frame.Model     `xorm:"-"`
}

func (this *Component) NewLog() *AccountLogModel {
	logEntity := this.OrmModel(&AccountLogModel{}).(*AccountLogModel)
	logEntity.SetTableName(this.group + "_account_log")
	return logEntity
}

// 指定主键
func (this *AccountLogModel) PrimaryKey() interface{} {
	return this.Id
}
