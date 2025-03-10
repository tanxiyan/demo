package captcha

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"time"

	"gitee.com/go-mao/mao/libs/try"

	"gitee.com/go-mao/mao/frame"
	"github.com/dchest/captcha"
)

var cacheStore = captcha.NewMemoryStore(100, time.Minute*10)

type Component struct {
	*frame.Taskline
	enable bool
}

func New(t *frame.Taskline, enable bool) *Component {
	object := new(Component)
	object.Taskline = t
	object.enable = enable
	return object
}

// 创建验证码并输出
func (this *Component) MakeImage() (sn string, base64Img string) {
	sn = captcha.NewLen(4)
	var data = bytes.NewBuffer(nil)
	if err := captcha.WriteImage(data, sn, 200, 80); err != nil {
		try.Throw(frame.CODE_WARN, "验证码图片创建失败：", err.Error())
	}
	imgByte, err := ioutil.ReadAll(data)
	if err != nil {
		try.Throw(frame.CODE_WARN, "验证码字节码生成失败：", err.Error())
	}
	base64Img = base64.StdEncoding.EncodeToString(imgByte)
	return
}

// 验证验证码（可验证两次）
func (this *Component) Validate(sn string, number string) {
	if !this.enable {
		return
	}
	digits := []byte(number)
	if bytes.Equal(cacheStore.Get(sn, true), digits) {
		return
	}
	if captcha.VerifyString(sn, number) {
		cacheStore.Set(sn, digits)
		return
	}
	try.Throw(frame.CODE_WARN, "验证码错误")
}
