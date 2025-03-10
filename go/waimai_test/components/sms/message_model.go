package sms

import (
	"time"

	"gitee.com/go-mao/mao/libs/utils"

	"gitee.com/go-mao/mao/frame"
)

const (
	STATE_SEND = 1 //已发送
	STATE_READ = 2 //已阅读
)

// 用户短信结构体
type MessageModel struct {
	Id          int64
	Index       string    `xorm:"varchar(40) index comment('索引')"`
	Code        string    `xorm:"varchar(12) comment('验证码')"`
	Phone       string    `xorm:"varchar(11) index comment('手机号码')"`
	State       int       `xorm:"tinyint(4) comment('状态')"`
	Mode        string    `xorm:"varchar(32) comment('消息类型')"`
	CreatedAt   time.Time `xorm:"created datetime comment('创建时间')"`
	UpdatedAt   time.Time `xorm:"updated datetime comment('更新时间')"`
	Version     int64     `xorm:"version"`
	debug       bool      `xorm:"-"` //调试模式不阅读
	frame.Model `xorm:"-"`
}

func (this *Component) NewMessage() *MessageModel {
	debug := this.Config().(*Config).isDebug()
	return this.OrmModel(&MessageModel{debug: debug}).(*MessageModel)
}

func (this *MessageModel) PrimaryKey() interface{} {
	return this.Id
}

// 表
func (this *MessageModel) TableName() string {
	return "SMS_Message"
}

func (this *MessageModel) makeIndex() string {
	return utils.Sha1(this.Phone, this.Code, this.Mode)
}

// 设为已阅读
func (this *MessageModel) SetRead() {
	if this.debug {
		return
	}
	this.State = STATE_READ
	this.Cols("State").Update()
}
