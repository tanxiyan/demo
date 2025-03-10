package user

type UserInterface interface {
	GetUserId() int64
	GetUserGroup() string
}

const (
	USER_STATE_ALLOW = 1
	USER_STATE_DENY  = -1
)

type LoginState int

const (
	LOGIN_STATE_OTHER  LoginState = -2 //异地登录
	LOGIN_STATE_NOPASS LoginState = -1 //验证失败
	LOGIN_STATE_PASS   LoginState = 1  //登录成功
)
