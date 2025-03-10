package permission

import (
	"time"

	"gitee.com/go-mao/mao/frame"
)

// 管理员角色
type UserRoleModel struct {
	Id          int64
	UserId      int64     `xorm:"int(11) index comment('管理员id')"`
	RoleId      int64     `xorm:"int(11)  comment('角色id')"`
	RoleName    string    `xorm:"-"`
	CreateTime  time.Time `xorm:"created   comment('创建时间')"`
	CreateBy    int64     `xorm:"bigint(20)"`
	frame.Model `xorm:"-"`
}

func (this *Component) NewUserRole() *UserRoleModel {
	userRoleEntity := this.OrmModel(&UserRoleModel{}).(*UserRoleModel)
	userRoleEntity.SetTableName(this.group + "_user_role")
	return userRoleEntity
}

func (this *UserRoleModel) PrimaryKey() interface{} {
	return this.Id
}
