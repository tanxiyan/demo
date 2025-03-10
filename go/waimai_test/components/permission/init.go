package permission

import (
	"gitee.com/go-mao/mao/frame"
)

type Component struct {
	group string
	*frame.Taskline
}

// prefix数据库表前缀名
func New(t *frame.Taskline, group string) *Component {
	object := new(Component)
	object.Taskline = t
	object.group = group
	return object
}

// 组件名称
func (this *Component) CompName() string {
	return this.group + "Permission"
}

func (this *Component) Models() []frame.ModelInterface {
	return []frame.ModelInterface{
		this.NewRole(),
		this.NewUserRole(),
	}
}

// 权限校验
func (this *Component) Check(user UserInterface, permissionId string) bool {
	if user.GetUserGroup() != this.group {
		return false
	}
	cache, ok := permissionCache.Load(user.GetUserId())
	if ok {
		if permissionsMap, ok := cache.(map[string]bool); ok {
			return permissionsMap[permissionId]
		}
	}
	permissions := this.RoleDomain().GetUserPermissons(user.GetUserId())
	permissionsMap := make(map[string]bool)
	for _, pemisId := range permissions {
		permissionsMap[pemisId] = true
	}
	permissionCache.Store(user.GetUserId(), permissionsMap)
	return permissionsMap[permissionId]
}
