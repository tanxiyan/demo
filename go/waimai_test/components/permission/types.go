package permission

type UserInterface interface {
	GetUserId() int64
	GetUserGroup() string
}

type UserGrouper interface {
	GetUserGroup() string
}
