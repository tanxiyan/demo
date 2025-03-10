package permission

import (
	"fmt"
	"sync"

	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
)

// 权限缓存，避免每次验证权限时读取数据，角色发生变化时自动清理缓存
var permissionCache = &sync.Map{}

// 角色服务
type roleDomain struct {
	*Component
}

func (this *Component) RoleDomain() *roleDomain {
	return &roleDomain{this}
}

// 角色列表
func (this *roleDomain) Select(listPtr interface{}) {
	db := this.OrmTable(this.NewRole())
	db.Select("*")
	if err := db.Find(listPtr); err != nil {
		try.Throw(frame.CODE_SQL, "查询角色错误：", err.Error())
	}
}

// 获取角色详情by Id
func (this *roleDomain) GetById(id int64) (*RoleModel, bool) {
	role := this.NewRole()
	ok := role.Match("Id", id).Get()
	return role, ok
}

// 获取角色详情by Id
func (this *roleDomain) GetByName(name string) (*RoleModel, bool) {
	role := this.NewRole()
	ok := role.Match("Name", name).Get()
	return role, ok
}

// 角色新增
func (this *roleDomain) Create(handler UserInterface, name, remark string, permissons []string) (r *RoleModel) {
	if _, ok := this.GetByName(name); ok {
		try.Throw(frame.CODE_WARN, "角色名称不可重复")
		return
	}
	roleEntity := this.NewRole()
	roleEntity.Name = name
	roleEntity.Remark = remark
	roleEntity.Permissions = permissons
	roleEntity.CreateBy = handler.GetUserId()
	if !roleEntity.Create() {
		try.Throw(frame.CODE_SQL, "角色添加失败")
	}
	//清空权限缓存
	permissionCache = &sync.Map{}
	return roleEntity
}

// 角色修改
func (this *roleDomain) Update(handler UserInterface, data RoleModel, cols ...string) {
	roleEntity, ok := this.GetById(data.Id)
	if !ok {
		try.Throw(frame.CODE_WARN, "该角色无效")
	}
	if r, ok := this.GetByName(data.Name); ok && r.Id != roleEntity.Id {
		try.Throw(frame.CODE_WARN, "角色名称不可重复")
		return
	}
	roleEntity.Name = data.Name
	roleEntity.Remark = data.Remark
	roleEntity.Permissions = data.Permissions
	roleEntity.UpdateBy = handler.GetUserId()
	if !roleEntity.Cols(cols...).Update() {
		try.Throw(frame.CODE_SQL, "角色修改失败")
	}
	//清空权限缓存
	permissionCache = &sync.Map{}
}

// 删除角色
func (this *roleDomain) Delete(roleIds []int64) {
	if total, _ := this.OrmTable(this.NewUserRole()).In("RoleId", roleIds).Count(); total > 0 {
		try.Throwf(frame.CODE_WARN, "该角色正在使用")
	}
	db := this.OrmTable(this.NewRole())
	_, err := db.In("Id", roleIds).Delete()
	if err != nil {
		try.Throwf(frame.CODE_WARN, "角色删除失败：%s", err.Error())
	}
	//清空权限缓存
	permissionCache = &sync.Map{}
}

// 给管理员角色授权,可授予多个角色,先删除之前的角色
func (this *roleDomain) GrantUser(handler, toUser UserInterface, roleIds []int64) {
	if total, _ := this.OrmSession().Table(this.NewRole()).In("Id", roleIds).Count(); len(roleIds) != int(total) {
		try.Throwf(frame.CODE_WARN, "角色异常")
	}
	this.Transaction(func() {
		//更新角色缓存
		permissionCache.Delete(toUser.GetUserId())
		table := this.NewUserRole()
		_, err := this.OrmSession().Where("UserId=?", toUser.GetUserId()).Delete(table)
		if err != nil {
			try.Throwf(frame.CODE_SQL, "用户%d授权删除旧角色失败：%s", toUser.GetUserId(), err.Error())
		}
		roles := make([]*UserRoleModel, 0)
		for _, roleId := range roleIds {
			roles = append(roles, &UserRoleModel{
				UserId:   toUser.GetUserId(),
				RoleId:   roleId,
				CreateBy: handler.GetUserId(),
			})
		}
		if _, err := this.OrmTable(table).InsertMulti(&roles); err != nil {
			try.Throwf(frame.CODE_SQL, "用户%d授权添加角色失败：%s", toUser.GetUserId(), err.Error())
		}
	})
}

// 获取管理员角色
func (this *roleDomain) GetUserRoles(userId int64, listPtr interface{}) {
	db := this.OrmTable(this.NewUserRole()).Alias("m")
	db.Join("INNER", this.NewRole().Alias("r"), "`r`.`Id`=`m`.`RoleId`")
	db.Where("`m`.`UserId`=?", userId)
	db.Select("`r`.`Name`,`r`.`Id`")
	if err := db.Find(listPtr); err != nil {
		try.Throwf(frame.CODE_SQL, "获取用户%d角色失败：%s", userId, err.Error())
	}
}

// 删除管理员某个角色
func (this *roleDomain) DeleteUserRoles(userIds []int64) {
	db := this.OrmSession().In("UserId", userIds)
	_, err := db.Delete(this.NewUserRole())
	if err != nil {
		try.Throwf(frame.CODE_SQL, "用户%d角色修改失败：%s", fmt.Sprint(userIds), err.Error())
	}
	//更新角色缓存
	for _, userId := range userIds {
		permissionCache.Delete(userId)
	}
}

// 获取管理员权限列表
func (this *roleDomain) GetUserPermissons(userId int64) []string {
	db := this.OrmTable(this.NewUserRole()).Alias("u")
	db.Join("INNER", this.NewRole().Alias("r"), "`r`.`Id`=`u`.`RoleId`")
	db.Where("`u`.`UserId`=?", userId)
	db.Select("`r`.`Permissions`")
	list := make([]struct {
		Permissions []string
	}, 0)
	if err := db.Find(&list); err != nil {
		try.Throwf(frame.CODE_SQL, "获取用户%d权限失败：%s", userId, err.Error())
	}
	permissons := make([]string, 0)
	for _, item := range list {
		permissons = append(permissons, item.Permissions...)
	}
	return permissons
}
