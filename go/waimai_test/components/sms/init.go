package sms

import (
	"fmt"
	"math/rand"
	"time"

	"gitee.com/go-mao/mao/libs/binding"
	"gitee.com/go-mao/mao/libs/try"
	"gitee.com/go-mao/mao/libs/utils"

	"gitee.com/go-mao/mao/frame"
	"github.com/northbright/aliyun/message"
)

type Component struct {
	*frame.Taskline
}

func New(t *frame.Taskline) *Component {
	object := new(Component)
	object.Taskline = t
	return object
}
func (this *Component) CompName() string {
	return "验证码"
}

// 配置
func (this *Component) Config() frame.ConfigInterface {
	return this.RenderConfig(&Config{})
}

func (this *Component) Models() []frame.ModelInterface {
	return []frame.ModelInterface{
		this.NewMessage(),
	}
}

// 发送验证码
// mobile=手机号码，mode=消息类型（用户自定义）
func (this *Component) Send(mobile string, mode string) {
	if !binding.IsMobile(mobile) {
		try.Throw(frame.CODE_WARN, "请输入正确的手机号码")
	}
	config := this.Config().(*Config)
	//选择接入商
	var sender Sender
	switch config.UseMode {
	case "ali":
		sender = this.aliDomain()
	case "chuanglan":
		sender = this.chuangLanDomain()
	case "market":
		sender = this.marketDomain()
	default:
		return
	}
	//单次发送时间限制
	msgInfo := this.NewMessage()
	lastSendAt := time.Now().Add(0 - time.Second*send_wait_time).Format(utils.FORMAT_DATE_TIME)
	if l, _ := this.OrmTable(msgInfo).Where("`Phone`=? and `CreatedAt`>=?", mobile, lastSendAt).Count(); l > 0 {
		try.Throw(frame.CODE_WARN, "请稍候再试")
	}
	//1小时发送限制
	lastSendAt = time.Now().Add(0 - time.Minute*60).Format(utils.FORMAT_DATE_TIME)
	if l, _ := this.OrmTable(msgInfo).Where("`Phone`=? and `CreatedAt`>=?", mobile, lastSendAt).Count(); l >= config.SendLimit {
		try.Throw(frame.CODE_WARN, fmt.Sprintf("手机号码%s已触发流控限制，请等待1小时候再试", mobile))
	}
	//发送验证码
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprint(rand.Int31(), "880321")[:config.CodeLen]
	this.Transaction(func() {
		sender.Send(mode, mobile, code)
		msgInfo.Phone = mobile
		msgInfo.Code = code
		msgInfo.State = STATE_SEND
		msgInfo.Mode = mode
		msgInfo.Index = msgInfo.makeIndex()
		msgInfo.Create()
	})
}

// 号码验证
func (this *Component) Validate(mobile, code, mode string) (*MessageModel, bool) {
	msgInfo := this.NewMessage()
	config := this.Config().(*Config)
	if config.isDebug() {
		return msgInfo, true
	}
	msgInfo.Phone = mobile
	msgInfo.Code = code
	msgInfo.Mode = mode
	sendTime := time.Now().Add(0 - time.Minute*time.Duration(config.Expire)).Format(utils.FORMAT_DATE_TIME)
	ok := msgInfo.Where("`Index`=? and `State`=? and `CreatedAt`>=?", msgInfo.makeIndex(), STATE_SEND, sendTime).Get()
	if !ok {
		return nil, false
	}
	if config.UseMode == "none" {
		msgInfo.debug = true
	}
	return msgInfo, true
}

// 查询参数
type FindArgs struct {
	Phone   string
	BeginAt string
	EndAt   string
}

// 查询验证码
func (this *Component) Find(args *FindArgs, page, limit int, listPtr interface{}) int64 {
	db := this.OrmTable(this.NewMessage())
	if args.Phone != "" {
		db.Where("`Mobile`=?", args.Phone)
	}
	if args.BeginAt != "" {
		db.And("`CreatedAt`>?", args.BeginAt)
	}
	if args.EndAt != "" {
		db.And("`CreatedAt`<?", args.EndAt)
	}
	db.Desc("Id")
	return this.FindPage(db, listPtr, page, limit)
}

// 阿里发送短信
func (this *Component) AliSendMessage(mobile string, templateCode, templateParam string) {
	if templateCode == "" {
		return
	}
	config := this.Config().(*Config)
	conf := config.Ali
	sms := message.NewClient(conf.AccessKeyID, conf.AccessKeySecret)
	ok, rsp, err := sms.SendSMS([]string{mobile}, conf.SignName, templateCode, templateParam)
	if !ok && rsp != nil {
		try.Throwf(frame.CODE_WARN, "短信发送失败", "阿里，发送手机：%s，内容：%s，错误信息：%s ", mobile, templateParam, rsp.Message)
	}
	if err != nil {
		try.Throwf(frame.CODE_WARN, "短信发送失败", "阿里，发送手机：%s，内容：%s，错误信息：%s ", mobile, templateParam, err.Error())
	}
}
