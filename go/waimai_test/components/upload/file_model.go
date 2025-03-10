package upload

import (
	"time"

	"gitee.com/go-mao/mao/frame"
)

// 上传文件模型（所有上传文件保存起来）
type FileModel struct {
	Id          int64
	UserId      int64     `xorm:"int(11) index comment('用户id')"`
	UserGroup   string    `xorm:"varchar(64) comment('用户所在组')"`
	Storage     string    `xorm:"varchar(32)  comment('存储类型')"`
	Mode        string    `xorm:"varchar(64) comment('上传类型')"`
	Src         string    `xorm:"varchar(255) comment('本地存储地址，方便删除')"`
	Url         string    `xorm:"varchar(255) comment('对外地址')"`
	Size        int64     `xorm:"int(11) comment('文件,kb')"`
	State       int       `xorm:"tinyint comment('资源状态')"`
	Md5         string    `xorm:"varchar(32) index comment('文件md5值')"`
	CreateTime  time.Time `xorm:"created index  comment('创建时间')"`
	frame.Model `xorm:"-"`
}

func (this *Component) NewFile() *FileModel {
	object := this.OrmModel(&FileModel{}).(*FileModel)
	return object
}

func (this *FileModel) TableName() string {
	return "upload_files"
}

func (this *FileModel) PrimaryKey() interface{} {
	return this.Id
}
