package wechat

import (
	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
	"github.com/ArtisanCloud/PowerWeChat/v2/src/miniProgram"
	"github.com/ArtisanCloud/PowerWeChat/v2/src/miniProgram/auth/response"
)

type Component struct {
	*frame.Taskline
}

func New(t *frame.Taskline) *Component {
	object := new(Component)
	object.Taskline = t
	return object
}

func (this *Component) Models() []frame.ModelInterface {
	return []frame.ModelInterface{}
}

func (this *Component) Config() frame.ConfigInterface {
	return this.RenderConfig(&Config{})
}

func (this *Component) CompName() string {
	return "base.wechat"
}

func (this *Component) GetConfig() *Config {
	return this.Config().(*Config)
}

func (this *Component) NewClient() *miniProgram.MiniProgram {
	config := this.GetConfig()
	miniProgramApp, err := miniProgram.NewMiniProgram(&miniProgram.UserConfig{
		AppID:     config.AppId,     // 小程序appid
		Secret:    config.AppSecret, // 小程序app secret
		HttpDebug: true,
		Log: miniProgram.Log{
			Level: "debug",
			File:  "./runtime/wechat.log",
		},
	})
	if err != nil {
		try.Throw(frame.CODE_FATAL, "Wechat初始化失败")
	}
	return miniProgramApp
}

// Login 小程序登录
func (this *Component) GetAccountInfo(code, phoneCode string) (res *response.ResponseCode2Session, phone string, err error) {
	client := this.NewClient()
	var sessionErr error
	res, sessionErr = client.Auth.Session(code)
	if sessionErr != nil {
		err = sessionErr
		return
	}
	phoneRes, getPhoneErr := client.PhoneNumber.GetUserPhoneNumber(phoneCode)
	if getPhoneErr != nil {
		err = getPhoneErr
		return
	}
	phone = phoneRes.PhoneInfo.PhoneNumber
	return
}

// Login 小程序登录
func (this *Component) GetOpenId(code string) (res *response.ResponseCode2Session, err error) {
	client := this.NewClient()
	var sessionErr error
	res, sessionErr = client.Auth.Session(code)
	if sessionErr != nil {
		err = sessionErr
		return
	}
	return
}
