package user

import (
	"time"

	"gitee.com/go-mao/mao/frame"
)

type UserModel struct {
	Id            int64
	Username      string    `xorm:"varchar(40) unique comment('登录账号')"`
	Phone         string    `xorm:"varchar(11) unique comment('手机号')"`
	Password      string    `xorm:"varchar(40) comment('登录密码sha1')"`
	Nickname      string    `xorm:"varchar(32) comment('昵称')"`
	Realname      string    `xorm:"varchar(32) comment('真实姓名')"`
	Remarks       string    `xorm:"varchar(255) comment('备注')"`
	Sex           string    `xorm:"varchar(1) comment('性别')"`
	Avatar        string    `xorm:"varchar(255) comment('头像地址')"`
	Email         string    `xorm:"varchar(255) comment('邮箱地址')"`
	State         int       `xorm:"tinyint(4) comment('状态')"`
	LoginErrNum   int       `xorm:"tinyint(4) comment('错误登录次数')"`
	CreateIP      string    `xorm:"varchar(64) comment('创建IP')"`
	CreateTime    time.Time `xorm:"created  datetime comment('创建时间')"`
	CreateBy      int64     `xorm:"bigint(20) comment('创建人')"`
	OpenId        string    `xorm:"varchar(32) comment('微信openId')"`
	LastLoginTime time.Time `xorm:"datetime comment('最后登录时间')"`
	UpdateBy      int64     `xorm:"bigint(20) comment('更新人')"`
	UpdateTime    time.Time `xorm:"updated datetime comment('最后修改时间')"`
	group         string    `xorm:"-"`
	frame.Model   `xorm:"-"`
}

// 实例化user model
func (this *Component) NewUser() *UserModel {
	userEntity := this.OrmModel(&UserModel{}).(*UserModel)
	userEntity.SetTableName(this.CompName())
	userEntity.group = this.group
	return userEntity
}

func (this UserModel) GetUserId() int64 {
	return this.Id
}

func (this UserModel) GetUserGroup() string {
	return this.group
}

// 设置主键值（
func (this *UserModel) PrimaryKey() interface{} {
	return this.Id
}
