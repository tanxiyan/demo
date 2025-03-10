package user

import (
	"time"

	"github.com/gin-gonic/gin"

	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
	"gitee.com/go-mao/mao/libs/utils"
)

type Component struct {
	group  string
	config *Config
	*frame.Taskline
}

func New(t *frame.Taskline, group string, config *Config) *Component {
	object := new(Component)
	object.Taskline = t
	object.config = config
	object.group = group
	return object
}

// 组件名称
func (this *Component) CompName() string {
	return this.group
}

// 数据库模型
func (this *Component) Models() []frame.ModelInterface {
	return []frame.ModelInterface{
		this.NewUser(),
		this.NewUserSession(),
	}
}

// 登录设备验证
func (this *Component) CheckDevice(device string) {
	if this.config.GetDevice(device) == nil {
		try.Throw(frame.CODE_WARN, "该设备不允许登录")
	}
}

// 获取uid，保证uid是从10000开始的
func (this *Component) getUid() int64 {
	const START_UID = 10000
	userEntity := this.NewUser()
	this.OrmSession().Cols("Id").Desc("Id").Get(userEntity)
	if userEntity.Id > 0 || userEntity.Id >= START_UID {
		return 0
	}
	return START_UID
}

// 创建用户
func (this *Component) Create(handler UserInterface, data *UserModel, clientIP string) *UserModel {
	if data.Username == "" && data.Phone == "" {
		try.Throw(frame.CODE_WARN, "用户名或手机号码不能为空")
	}
	userEntity := this.NewUser()
	userEntity.Id = this.getUid()
	if data.Username != "" {
		userEntity.Username = data.Username
		userEntity.Cols("Username")
	}
	if data.Phone != "" {
		userEntity.Phone = data.Phone
		userEntity.Cols("Phone")
	}
	if data.Email != "" {
		userEntity.Email = data.Email
		userEntity.Cols("Email")
	}
	if data.Nickname != "" {
		userEntity.Nickname = data.Nickname
		userEntity.Cols("Nickname")
	}
	if data.Avatar != "" {
		userEntity.Avatar = data.Avatar
		userEntity.Cols("Avatar")
	}
	if data.State != 0 {
		userEntity.State = data.State
		userEntity.Cols("State")
	}
	if data.Realname != "" {
		userEntity.Realname = data.Realname
		userEntity.Cols("Realname")
	}
	if data.Remarks != "" {
		userEntity.Remarks = data.Remarks
		userEntity.Cols("Remarks")
	}
	userEntity.CreateIP = clientIP
	userEntity.Password = utils.Sha1(data.Password)
	userEntity.Cols("CreateIP", "Password")
	if handler != nil {
		userEntity.CreateBy = handler.GetUserId()
		userEntity.Cols("CreateBy")
	}
	this.validate(0, userEntity)
	if !userEntity.Cols("Id", "Nickname", "Password", "State", "CreateBy", "CreateTime", "CreateIP", "UpdateTime").Create() {
		try.Throw(frame.CODE_SQL, "用户创建失败")
	}
	return userEntity
}

// 修改
func (this *Component) Update(handler UserInterface, userId int64, data *UserModel) {
	userEntity := this.GetByUserId(userId)
	if userEntity == nil {
		try.Throwf(frame.CODE_WARN, "用户%d无效", userId)
	}
	data.Id = userId
	if data.Username != "" {
		userEntity.Username = data.Username
		userEntity.Cols("Username")
	}
	if data.Phone != "" {
		userEntity.Phone = data.Phone
		userEntity.Cols("Phone")
	}
	if data.Email != "" {
		userEntity.Email = data.Email
		userEntity.Cols("Email")
	}
	if data.Password != "" {
		userEntity.Password = utils.Sha1(data.Password)
		userEntity.Cols("Password")
	}
	if data.Nickname != "" {
		userEntity.Nickname = data.Nickname
		userEntity.Cols("Nickname")
	}
	if data.Avatar != "" {
		userEntity.Avatar = data.Avatar
		userEntity.Cols("Avatar")
	}
	if data.State != 0 {
		userEntity.State = data.State
		userEntity.Cols("State")
	}
	if data.Realname != "" {
		userEntity.Realname = data.Realname
		userEntity.Cols("Realname")
	}
	if data.Remarks != "" {
		userEntity.Remarks = data.Remarks
		userEntity.Cols("Remarks")
	}
	if handler != nil {
		userEntity.UpdateBy = handler.GetUserId()
		userEntity.Cols("UpdateBy")
	}
	this.validate(userId, userEntity)
	this.Transaction(func() {
		if !userEntity.Update() {
			try.Throwf(frame.CODE_SQL, "用户%d修改失败", userId)
		}
	})
}

// 数据验证
func (this *Component) validate(userId int64, data *UserModel) {
	if data.Username != "" {
		if this.IsUsernameExists(data.Username, userId) {
			try.Throwf(frame.CODE_WARN, "用户名%s不可用", data.Username)
		}
	}
	if data.Phone != "" {
		if this.IsPhoneExists(data.Phone, userId) {
			try.Throwf(frame.CODE_WARN, "手机号码%s不可用", data.Phone)
		}
	}
	if data.Email != "" {
		if this.IsEmailExists(data.Email, userId) {
			try.Throwf(frame.CODE_WARN, "邮箱%s不可用", data.Email)
		}
	}
	if data.State > 0 {
		if utils.ArrayNotIn(data.State, USER_STATE_ALLOW, USER_STATE_DENY) {
			try.Throw(frame.CODE_WARN, "用户状态错误")
		}
	}
}

// 验证用户名是否有效
func (this *Component) IsUsernameExists(username string, omitUserId int64) bool {
	userEntity := this.NewUser()
	userEntity.Where("`Username`=? and `Id`!=?", username, omitUserId)
	return userEntity.Exists()
}

// 验证手机是否有效
func (this *Component) IsPhoneExists(phone string, omitUserId int64) bool {
	userEntity := this.NewUser()
	userEntity.Where("`Phone`=? and `Id`!=?", phone, omitUserId)
	return userEntity.Exists()
}

// 验证邮箱是否有效
func (this *Component) IsEmailExists(email string, omitUserId int64) bool {
	userEntity := this.NewUser()
	userEntity.Where("`Email`=? and `Id`!=?", email, omitUserId)
	return userEntity.Exists()
}

// 用户名登录
func (this *Component) LoginByUsername(username, password string) (*UserModel, bool) {
	userEntity := this.NewUser()
	ok := userEntity.Match("Username", username).Get()
	if ok && utils.Sha1(password) == userEntity.Password {
		return userEntity, true
	}
	return userEntity, false
}

// 手机号码登录
func (this *Component) LoginByPhone(phone, password string) (*UserModel, bool) {
	userEntity := this.NewUser()
	ok := userEntity.Match("Phone", phone).Get()
	if ok && utils.Sha1(password) == userEntity.Password {
		return userEntity, true
	}
	return userEntity, false
}

// 验证登录设备，并创建登录token，登录成功后解除锁定
func (this *Component) SaveLogin(user UserInterface, ip, device string) *UserSessionModel {
	deviceEntity := this.config.GetDevice(device)
	if deviceEntity == nil {
		try.Throw(frame.CODE_WARN, "该设备不允许登录")
	}
	loginEntity := this.NewUserSession()
	loginEntity.UserId = user.GetUserId()
	loginEntity.Device = device
	loginEntity.IP = ip
	loginEntity.Token = utils.Sha1(time.Now().UnixNano(), user.GetUserId(), device)
	loginEntity.Expired = time.Now().Add(time.Minute * time.Duration(deviceEntity.MaxSession))
	userEntity := this.GetByUserId(user.GetUserId(), "Id", "State")
	if userEntity == nil {
		try.Throwf(frame.CODE_WARN, "保存用户登录信息失败，用户%d无效", user.GetUserId())
	}
	this.Transaction(func() {
		userEntity.State = USER_STATE_ALLOW
		userEntity.LoginErrNum = 0
		userEntity.LastLoginTime = time.Now()
		if !userEntity.Cols("State", "LoginErrNum", "LastLoginTime").Must("LoginErrNum").Update() {
			try.Throw(frame.CODE_SQL, "更新用户状态错误，用户", user.GetUserId())
		}
		this.OrmTable(this.NewUserSession().TableName()).Where("`UserId`=? and `Expired`>NOW() and `Device`=?", userEntity.Id, device).Update(gin.H{
			"Expired": "0001-01-01 00:00:00",
		})
		if !loginEntity.Create() {
			try.Throw(frame.CODE_SQL, "登录信息保存失败，用户", user.GetUserId())
		}
	})
	return loginEntity
}

// 根据token验证登录信息，每隔10分钟动延长到期时间
// return 登录对象，是否登录成功，是否异地
func (this *Component) CheckLogin(token string, device string) (*UserSessionModel, LoginState) {
	deviceEntity := this.config.GetDevice(device)
	if deviceEntity == nil {
		return nil, LOGIN_STATE_NOPASS
	}
	currentTime := time.Now()
	loginEntity := this.NewUserSession()
	db := this.OrmTable(loginEntity).Where("`Token`=? and `Device`=?", token, device)
	ok, err := db.Desc("Id").Get(loginEntity)
	if err != nil {
		try.Throw(frame.CODE_SQL, "登录验证异常：", err.Error())
	}
	if loginEntity.Expired.IsZero() {
		return nil, LOGIN_STATE_OTHER
	}
	if loginEntity.Expired.Before(currentTime) {
		return nil, LOGIN_STATE_NOPASS
	}
	if !ok {
		return nil, LOGIN_STATE_NOPASS
	}
	minutes := deviceEntity.MaxSession / 2
	if loginEntity.UpdateTime.Add(time.Minute * time.Duration(minutes)).Before(currentTime) {
		loginEntity.Expired = time.Now().Add(time.Minute * time.Duration(deviceEntity.MaxSession))
		loginEntity.Cols("Expired").Update()
	}
	return loginEntity, LOGIN_STATE_PASS
}

// 增加错误登录次数
func (this *Component) AddLoginErrNum(userId int64) {
	userEntity := this.GetByUserId(userId, "Id", "LoginErrNum")
	if userEntity == nil {
		return
	}
	userEntity.LoginErrNum += 1
	if userEntity.LoginErrNum > this.config.MaxErrLogin {
		userEntity.State = USER_STATE_DENY
		userEntity.Cols("State")
	}
	if !userEntity.Cols("LoginErrNum").Update() {
		try.Throw(frame.CODE_WARN, " 无法更新用户", userId, "错误登录信息")
	}
}

// 获取用户信息
func (this *Component) GetByUserId(userId int64, cols ...string) *UserModel {
	userEntity := this.NewUser()
	if userEntity.Where("`Id`=?", userId).Cols(cols...).Get() {
		return userEntity
	}
	return nil
}

// 根据用户名查询用户对象
func (this *Component) GetByUsername(username string, cols ...string) (*UserModel, bool) {
	userEntity := this.NewUser()
	ok := userEntity.Match("Username", username).Cols(cols...).Get()
	return userEntity, ok
}

// 根据手机号码查询用户对象
func (this *Component) GetByPhone(phone string) (*UserModel, bool) {
	userEntity := this.NewUser()
	ok := userEntity.Match("Phone", phone).Get()
	return userEntity, ok
}

// 根据邮箱地址查询用户对象
func (this *Component) GetByEmail(email string) (*UserModel, bool) {
	userEntity := this.NewUser()
	ok := userEntity.Match("Email", email).Get()
	return userEntity, ok
}

// 获取用户账户信息
func (this *Component) GetByAccount(usernameOrPhone string) (*UserModel, bool) {
	userEntity := this.NewUser()
	ok := userEntity.Where("`Username`=? or `Phone`=?", usernameOrPhone, usernameOrPhone).Get()
	return userEntity, ok
}

// 批量删除
func (this *Component) Delete(ids []int64) {
	this.Transaction(func() {
		db := this.OrmTable(this.NewUser())
		if _, err := db.In("Id", ids).Delete(); err != nil {
			try.Throw(frame.CODE_SQL, "用户删除失败", ids)
		}
		db2 := this.OrmTable(this.NewUserSession())
		if _, err := db2.In("UserId", ids).Delete(); err != nil {
			try.Throw(frame.CODE_SQL, "用户登录信息删除失败", ids)
		}
	})
}

// 更新登录的token
func (this *Component) Logout(token string, device string) {
	loginEntity := this.NewUserSession()
	db := this.OrmTable(loginEntity).Where("`Token`=? and `Device`=?", token, device)
	_, err := db.Desc("Id").Get(loginEntity)
	if err != nil {
		try.Throw(frame.CODE_SQL, "登录验证异常：", err.Error())
	}
	loginEntity.Expired = time.Now().Add(0 - time.Hour*1)
	if !loginEntity.Cols("Expired").Update() {
		try.Throwf(frame.CODE_SQL, "用户%d退出登录失败", loginEntity.UserId)
	}
}

// 查询登陆日志
func (this *Component) FindLoginSessions(userId int64, page, limit int, listPtr interface{}) int64 {
	db := this.OrmTable(this.NewUserSession()).Alias("l")
	db.Join("INNER", this.NewUser().Alias("u"), "`l`.`UserId`=`u`.`Id`")
	if userId > 0 {
		db.Where("`l`.`UserId`=?", userId)
	}
	db.Desc("Id")
	db.Select("`l`.*,`u`.`Username`")
	return this.FindPage(db, listPtr, page, limit)
}
