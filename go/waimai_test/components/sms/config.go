package sms

import (
	"errors"
	"strings"

	"gitee.com/go-mao/mao/frame"
)

// 组件配置
type Config struct {
	UseMode   string `name:"接入商" mode:"single" desc:"" options:"none:关闭 market:服务市场 ali:阿里"` // chuanglan:创蓝
	CodeLen   int    `name:"验证码位数" mode:"single" desc:"" options:"4:4位 6:6位"`
	Expire    int    `name:"验证码有效时长" desc:"单位：分，发送后至该时间范围内有效" valid:"required,gt=0,lt=1000->验证码有效期为1-1000分钟以内"`
	SendLimit int64  `name:"一小时发送次数限制"  valid:"required,gt=0,lt=1000->一小时发送次数限制为1-1000之间"`
	Market    struct {
		SignName string `name:"短信签名名称" desc:"自定义配置，可以是你网站的名字，请不要超过6个字，不要使用符号"`
		Remarks  string `name:"使用说明" mode:"html"`
	} `name:"服务市场" desc:"请到服务市场购买额度后使用"`
	Ali struct {
		AccessKeyID     string `name:"接入ID"`
		AccessKeySecret string `name:"接入密匙"`
		SignName        string `name:"签名"`
		TemplateCode    string `name:"模板"`
	} `name:"阿里短信" desc:"接入商配置"`
	ChuangLan struct {
		Account  string `name:"账号"`
		Password string `name:"密码"`
		Url      string `name:"接口地址"`
		Template string `name:"模板"`
	} `name:"-" desc:"接入商配置"`
}

// 默认配置
func (this *Config) Default() frame.ConfigInterface {
	this.UseMode = "market"
	this.Expire = 5    //默认5分钟有效
	this.SendLimit = 5 //默认5次
	this.Market.SignName = "外卖提示"
	this.Market.Remarks = "请到服务市场购买短信额度"
	return this
}

// 配置名称
func (this *Config) ConfigName() string {
	return "短信验证码"
}

// 配置别名
func (this *Config) ConfigAlias() string {
	return "sms"
}

// 获取钱
func (this *Config) BeforeGet() {
}

func (this *Config) Validate() error {
	if strings.Index(this.Market.SignName, "【") > -1 || strings.Index(this.Market.SignName, "】") > -1 {
		return errors.New("签名不能出现字符【】")
	}
	if len([]rune(this.Market.SignName)) > 6 {
		return errors.New("签名字符不能超过6位")
	}
	return nil
}

// 是否未调试模式
func (this *Config) isDebug() bool {
	return this.UseMode == "none"
}
