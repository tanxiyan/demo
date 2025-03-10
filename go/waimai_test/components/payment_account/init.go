package payment_account

import (
	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
	"strings"
	"time"
)

type Component struct {
	prefix string
	*frame.Taskline
	config *Config
}

// 模块实例化唯一入口
func New(t *frame.Taskline, config Config) *Component {
	object := new(Component)
	object.Taskline = t
	object.config = &config
	return object
}

func (this *Component) Models() []frame.ModelInterface {
	return []frame.ModelInterface{
		this.NewAccount(),
	}
}

func (this *Component) CompName() string {
	return this.prefix + "_Payment_Account"
}

// 查询用户支付账号
// user 用户
// accountId 支付帐户id
// 返回：支付帐户，数据是否存在
func (this *Component) GetAccount(user UserInterface, accountId int64) (*AccountModel, bool) {
	account := this.NewAccount()
	ok := account.Where("`UserId`=? and `Id`=?", user.GetUserId(), accountId).Get()
	if !ok {
		return nil, false
	}
	return account, true
}

// 用户帐户查询
// user 用户信息
// page 页
// limit 每页数量
func (this *Component) Find(userId int64, account, realname string, page, limit int, listPtr interface{}) int64 {
	sx := this.OrmTable(this.NewAccount())
	if userId > 0 {
		sx.Where("`UserId`=?", userId)
	}
	if account != "" {
		sx.Where("`Account` like ?", "%"+account+"%")
	}
	if realname != "" {
		sx.Where("`Realname` like ?", "%"+realname+"%")
	}
	sx.Desc("`Id`")
	return this.FindPage(sx, listPtr, page, limit)
}

// 用户帐户列表
// user 用户信息
func (this *Component) FindAll(userId int64, listPtr interface{}) {
	sx := this.OrmTable(this.NewAccount())
	sx.Where("`UserId`=?", userId)
	err := sx.Find(listPtr)
	if err != nil {
		try.Throwf(frame.CODE_SQL, "用户%d账户列表查询失败", userId)
	}
}

// 账号绑定,如果不允许重复绑定的则删除之前绑定的数据
// user 用户
// nickname 昵称
// headimage 头像
// accountNO 账号
// paymentTypeName 支付方式名称
// 返回：错误
func (this *Component) BindAccount(user UserInterface, data AccountModel) (d *AccountModel) {
	if !this.config.CheckAccountType(data.AccountType) {
		try.Throw(frame.CODE_WARN, "账号类型错误")
		return
	}
	data.Account = strings.TrimSpace(data.Account)
	accountInfo := this.NewAccount()
	exists := accountInfo.Where("`UserId`=? and `Account`=? and `AccountType`=?", user.GetUserId(), data.Account, data.AccountType).Get()
	limitDay := this.config.AllowUpdateWaitDay
	if exists && limitDay > 0 && accountInfo.UpdatedTime.AddDate(0, 0, limitDay).After(time.Now()) {
		try.Throwf(frame.CODE_WARN, "账号绑定后%d天内不可修改", limitDay)
		return
	}
	accountInfo.UserId = user.GetUserId()
	accountInfo.AccountType = data.AccountType
	accountInfo.NickName = data.NickName
	accountInfo.HeadImage = data.HeadImage
	accountInfo.Account = data.Account
	accountInfo.Realname = data.Realname
	accountInfo.Qrcode = data.Qrcode
	var ok bool
	if exists {
		ok = accountInfo.Update()
	} else {
		ok = accountInfo.Create()
	}
	if !ok {
		try.Throw(frame.CODE_SQL, "账号绑定失败")
		return
	}
	return accountInfo
}

func (this *Component) DeleteByUser(user UserInterface, id int64) {
	account, ok := this.GetAccount(user, id)
	if !ok {
		try.Throwf(frame.CODE_SQL, "用户%d的账户%d查询失败", user.GetUserId(), id)
		return
	}
	limitDay := this.config.AllowUpdateWaitDay
	if limitDay > 0 && account.UpdatedTime.AddDate(0, 0, limitDay).After(time.Now()) {
		try.Throwf(frame.CODE_WARN, "账号【%s】绑定后%d天内不可删除", account.Account, limitDay)
		return
	}
	if !account.Delete() {
		try.Throwf(frame.CODE_SQL, "支付账户删除失败", user.GetUserId(), id)
	}
}
