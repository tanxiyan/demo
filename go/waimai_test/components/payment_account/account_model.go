package payment_account

import (
	"time"

	"gitee.com/go-mao/mao/frame"
)

// 支付账号
// 一个用户可以绑定多个支付帐号，当平台付款给用用户时使用
type AccountModel struct {
	Id          int64
	UserId      int64     `xorm:"int(11) index comment('用户id')"`
	AccountType string    `xorm:"varchar(32) comment('帐户类型')"`
	Realname    string    `xorm:"varchar(64) comment('真实姓名')"`
	NickName    string    `xorm:"varchar(64) comment('昵称')"`
	HeadImage   string    `xorm:"varchar(255) comment('头像')"`
	Qrcode      string    `xorm:"varchar(255) comment('收款二维码')"`
	Account     string    `xorm:"varchar(255) index comment('绑定的帐户id')"`
	CreateTime  time.Time `xorm:"created comment('创建时间')"`
	UpdatedTime time.Time `xorm:"updated comment('最后修改时间')"`
	frame.Model `xorm:"-"`
	group       string `xorm:"-"`
}

func (this *Component) NewAccount() *AccountModel {
	account := this.OrmModel(&AccountModel{group: this.config.Group}).(*AccountModel)
	account.SetTableName(this.config.Group + "_payment_account")
	return account
}

func (this *AccountModel) PrimaryKey() interface{} {
	return this.Id
}
