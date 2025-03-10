package payments

import (
	"bufio"
	"errors"
	"gitee.com/go-mao/mao/frame"
	"os"
)

// 配置
type Config struct {
	ApiDomainName string       `name:"通知域名" desc:"支付接口通知域名，填管理端域名即可，格式【http://test.com】"`
	Wechat        WechatConfig `name:"微信商户号"`
	M2MConfig     M2MConfig    `name:""`
	Alipay        AlipayConfig `name:""`
}

type M2MConfig struct {
	AppId      string `name:"AppId"`
	PublicKey  string `name:"公匙" mode:"text"`
	PrivateKey string `name:"私匙" mode:"text"`
}

// 支付宝配置
type AlipayConfig struct {
	AppId      string `name:"AppId"`
	PublicKey  string `name:"公匙" mode:"text"`
	PrivateKey string `name:"私匙" mode:"text"`
}

type WechatConfig struct {
	AppId       string `name:"小程序、公众号或者企业微信的appId"`
	MchId       string `name:"微信支付商户号ID"`
	MchApiV3Key string `name:"APIv3密钥"`
	Key         string `name:"APIv2密钥"`
	CertPath    string `name:"证书公钥(cert.pem)内容" mode:"text"`
	KeyPath     string `name:"证书私钥(key.pem)内容" mode:"text"`
	SerialNo    string `name:"证书序列号"`
}

// 默认配置
func (this *Config) Default() frame.ConfigInterface {
	if this.Wechat.CertPath == "" && this.Wechat.KeyPath != "" {
		this.Validate()
	}
	return this
}

// 配置名称
func (this *Config) ConfigName() string {
	return "收款配置"
}

// 配置别名
func (this *Config) ConfigAlias() string {
	return "base.payments"
}

func (this *Config) BeforeGet() {

}

func (this *Config) Validate() error {
	file, err := os.OpenFile(CertPathURL, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return errors.New("open CertPathURL err：" + err.Error())
	}
	defer file.Close()
	str := this.Wechat.CertPath
	writer := bufio.NewWriter(file)
	writer.WriteString(str)
	writer.Flush()

	file2, err := os.OpenFile(KeyPathURL, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return errors.New("open KeyPathURL err" + err.Error())
	}
	defer file2.Close()
	str2 := this.Wechat.KeyPath
	writer2 := bufio.NewWriter(file2)
	writer2.WriteString(str2)
	writer2.Flush()
	return nil
}
