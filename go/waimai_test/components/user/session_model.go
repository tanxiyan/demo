package user

import (
	"sync"
	"time"

	"gitee.com/go-mao/mao/frame"
)

/*
用户登录信息保存，类似于session
可以保存相关数据到会话中
*/
type UserSessionModel struct {
	Id          int64
	UserId      int64             `xorm:"int(11) index comment('用户id')"`
	IP          string            `xorm:"varchar(64) comment('登录ip')"`
	Device      string            `xorm:"varchar(10) index comment('登录终端类型')"`
	Token       string            `xorm:"varchar(40) index comment('登录口令sha1')"`
	Value       map[string]string `xorm:"text comment('存储数据')"`
	Expired     time.Time         `xorm:"datetime comment('到期时间')"`
	CreateTime  time.Time         `xorm:"created datetime comment('登录时间')"`
	UpdateTime  time.Time         `xorm:"updated datetime comment('修改时间')"`
	group       string            `xorm:"-"`
	frame.Model `xorm:"-"`
	sync.Mutex  `xorm:"-"`
}

// 实例化user login model
func (this *Component) NewUserSession() *UserSessionModel {
	loginEntity := this.OrmModel(&UserSessionModel{}).(*UserSessionModel)
	loginEntity.SetTableName(this.CompName() + "_session")
	loginEntity.Value = make(map[string]string)
	loginEntity.group = this.group
	return loginEntity
}

// 主键
func (this *UserSessionModel) PrimaryKey() any {
	return this.Id
}

// 操作者编号
func (this *UserSessionModel) GetUserId() int64 {
	return this.UserId
}

// 操作者用户名
func (this *UserSessionModel) GetUserGroup() string {
	return this.group
}

// 设置一个session值，必须要调用SaveVal才能保存
func (this *UserSessionModel) SetVal(key, value string) {
	this.Lock()
	defer this.Unlock()
	this.Value[key] = value
}

func (this *UserSessionModel) GetVal(key string) string {
	this.Lock()
	defer this.Unlock()
	val, ok := this.Value[key]
	if !ok {
		val = ""
	}
	return val
}

func (this *UserSessionModel) DelVal(key string) {
	this.Lock()
	defer this.Unlock()
	delete(this.Value, key)
}

func (this *UserSessionModel) ClearVal() {
	this.Lock()
	defer this.Unlock()
	this.Value = make(map[string]string)
}

// 保存
func (this *UserSessionModel) SaveVal() {
	this.Cols("Value").Update()
}
