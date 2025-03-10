package permission

import (
	"time"

	"gitee.com/go-mao/mao/frame"
)

// 角色
type RoleModel struct {
	Id          int64
	Name        string    `xorm:"unique varchar(32) comment('角色名称')"`
	Remark      string    `xorm:"varchar(255) comment('备注')"`
	Permissions []string  `xorm:"longtext json comment('权限')"`
	CreateBy    int64     `xorm:"bigint(20)  comment('创建人')"`
	CreateTime  time.Time `xorm:"created  comment('创建时间')"`
	UpdateBy    int64     `xorm:"bigint(20)  comment('创建人')"`
	UpdateTime  time.Time `xorm:"updated  comment('最后修改时间')"`
	frame.Model `xorm:"-"`
}

func (this *Component) NewRole() *RoleModel {
	roleEntity := this.OrmModel(&RoleModel{}).(*RoleModel)
	roleEntity.Permissions = make([]string, 0)
	roleEntity.SetTableName(this.group + "_role")
	return roleEntity
}

func (this *RoleModel) PrimaryKey() interface{} {
	return this.Id
}
