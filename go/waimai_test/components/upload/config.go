package upload

import "gitee.com/go-mao/mao/frame"

type Config struct {
	Storage     string      `name:"上传模式" mode:"single" desc:"" options:"local:本地 qiniu:七牛"`
	MaxSize     int64       `name:"最大上传大小/M" valid:"required,gt=0,lt=100->上传文件限制为1-100之间" desc:"单位：M"`
	AllowExts   string      `name:"可上传文件类型"  desc:"使用空格分开，例如.exe .jpg"`
	LocalConfig LocalConfig `name:"本地上传配置"`
	QiniuConfig QiniuConfig `name:"七牛上传配置"`
	group       string      `name:"-"`
}

type LocalConfig struct {
	Domain string `name:"访问路径" desc:"web访问的初始路由组，例如：/files,网页访问地址为http://demo.com/files/文件"`
	Dir    string `name:"上传到本地的目录" desc:"请输入程序所在的的相对目录"`
}

type QiniuConfig struct {
	Domain    string `name:"图片访问域名" valid:"url->域名格式为http://www.demo.com"`
	Bucket    string `name:"Bucket"`
	Zone      string `name:"机房" mode:"single" options:"hd:华东 hb:华北 hn:华南 as0:东南亚 na0:北美 cn-east-2:华东-浙江2"`
	AccessKey string `name:"AccessKey"`
	Secret    string `name:"Secret"`
}

// 默认配置
func (this *Config) Default() frame.ConfigInterface {
	this.AllowExts = ".png .jpg .jpeg"
	this.MaxSize = 5 //默认5M
	this.Storage = "local"
	this.LocalConfig = LocalConfig{
		Domain: "/files",
		Dir:    "./files",
	}
	return this
}

// 配置名称
func (this *Config) ConfigName() string {
	return "上传"
}

// 配置别名
func (this *Config) ConfigAlias() string {
	return "upload"
}

// 获取钱
func (this *Config) BeforeGet() {
}

// 验证
func (this *Config) Validate() error {
	this.LocalConfig = LocalConfig{
		Domain: "/files",
		Dir:    "./files",
	}
	return nil
}
