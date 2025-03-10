package account

type UserInterface interface {
	GetUserId() int64
	GetUserGroup() string
}
