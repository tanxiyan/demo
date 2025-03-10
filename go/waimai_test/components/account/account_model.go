package account

import (
	"time"

	"gitee.com/go-mao/mao/frame"
)

type AccountModel struct {
	Id          int64
	UserId      int64     `xorm:"notnull int(11) unique comment('用户id')"`
	Balance     int64     `xorm:"bigint(20) index comment('账户余额')"`
	CreateTime  time.Time `xorm:"created   comment('创建时间')"`
	UpdatedAt   time.Time `xorm:"updated   comment('最后修改时间')"`
	Version     int64     `xorm:"version"`
	frame.Model `xorm:"-"`
}

func (this *Component) NewAccount() *AccountModel {
	accountEntity := this.OrmModel(&AccountModel{}).(*AccountModel)
	accountEntity.SetTableName(this.group + "_account")
	return accountEntity
}

// 指定主键
func (this *AccountModel) PrimaryKey() interface{} {
	return this.Id
}
