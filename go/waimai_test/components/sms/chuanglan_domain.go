package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
)

// 创蓝短信验证码
type chuangLanDomain struct {
	*Component
}

func (this *Component) chuangLanDomain() *chuangLanDomain {
	return &chuangLanDomain{this}
}

// 创蓝短信发送
func (this *chuangLanDomain) Send(mode, phone string, code string) {
	conf := this.Config().(*Config).ChuangLan
	params := make(map[string]interface{})
	params["account"] = conf.Account   //创蓝API账号
	params["password"] = conf.Password //创蓝API密码
	params["phone"] = phone            //手机号码
	params["msg"] = url.QueryEscape(fmt.Sprintf(conf.Template, code))
	params["report"] = "true"
	bytesData, err := json.Marshal(params)
	if err != nil {
		try.Throw(frame.CODE_WARN, "验证码发送失败", "创蓝，params数据错误：%s", err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", conf.Url, reader)
	if err != nil {
		try.Throw(frame.CODE_WARN, "验证码发送失败", "创蓝，请求创建失败：%s", err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		try.Throw(frame.CODE_WARN, "验证码发送失败", "创蓝，请求失败：%s", err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		try.Throw(frame.CODE_WARN, "验证码发送失败", "创蓝，读取回执错误：%s", err.Error())
		return
	}
	defer resp.Body.Close()
	if bytes.Contains(respBytes, []byte("code")) {
		type Response struct {
			Code string `json:"code"`
		}
		response := &Response{}
		readErr := json.Unmarshal(respBytes, response)
		if readErr != nil {
			try.Throwf(frame.CODE_WARN, "验证码发送失败", "创蓝，返回数据json解析错误：%s，返回内容：%s", readErr.Error(), string(respBytes))
		}
		if response.Code != "0" {
			try.Throwf(frame.CODE_WARN, "验证码发送失败", "创蓝，错误：%s", string(respBytes))
		}
	}
}
