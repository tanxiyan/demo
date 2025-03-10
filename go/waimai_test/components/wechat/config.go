package wechat

import (
	"gitee.com/go-mao/mao/frame"
)

// 配置
type Config struct {
	AppId     string `name:"小程序、公众号或者企业微信的appId"`
	AppSecret string `name:"小程序、公众号或者企业微信的appSecret"`
}

// 默认配置
func (this *Config) Default() frame.ConfigInterface {
	return this
}

// 配置名称
func (this *Config) ConfigName() string {
	return "微信登录配置"
}

// 配置别名
func (this *Config) ConfigAlias() string {
	return "base.wechat"
}

func (this *Config) BeforeGet() {

}

func (this *Config) Validate() error {
	return nil
}
