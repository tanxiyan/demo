package upload

import (
	"bytes"
	"io"
	"mime/multipart"

	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
)

type Component struct {
	config *Config
	*frame.Taskline
}

func New(t *frame.Taskline) *Component {
	object := new(Component)
	object.Taskline = t
	return object
}

// 组件名称
func (this *Component) CompName() string {
	return "Upload"
}

// 数据库模型
func (this *Component) Models() []frame.ModelInterface {
	return []frame.ModelInterface{
		this.NewFile(),
	}
}

// 静态文件处理
func (this *Component) FileStatic(eg *frame.WebEngine) {
	config := this.Config().(*Config)
	if config.Storage == STORAGE_LOCAL {
		eg.Engine().Static(config.LocalConfig.Domain, config.LocalConfig.Dir)
	}
}

// 配置
func (this *Component) Config() frame.ConfigInterface {
	return this.RenderConfig(&Config{})
}

// 上传文件
func (this *Component) Upload(data io.Reader, filename string, size int64) *uploader {
	return this.newUploader(data, filename, size)
}

// 上传表单提交的文件
func (this *Component) UploadPostFile(upfile *multipart.FileHeader) *uploader {
	file, e := upfile.Open()
	if e != nil {
		try.Throwf(frame.CODE_WARN, "上传文件读取失败：", e.Error())
	}
	defer file.Close()
	return this.newUploader(file, upfile.Filename, upfile.Size)
}

// 上传文件
func (this *Component) UploadBytes(filename string, data []byte) *uploader {
	file := bytes.NewBuffer(data)
	return this.newUploader(file, filename, int64(len(data)))
}

// 批量删除资源文件
func (this *Component) Remove(owner UserInterface, ids []int64) {
	config := this.Config().(*Config)
	this.RemoveByConfig(owner, ids, config.LocalConfig, config.QiniuConfig)
}

// 批量删除资源文件
func (this *Component) RemoveByConfig(owner UserInterface, ids []int64, localConfig LocalConfig, qiniuConfig QiniuConfig) {
	upfiles := make([]*FileModel, 0)
	db := this.OrmTable(this.NewFile())
	if owner != nil {
		db.And("`UserId`=? and `UserGroup`=?", owner.GetUserId(), owner.GetUserGroup())
	}
	if err := db.In("Id", ids).Find(&upfiles); err != nil {
		try.Throw(frame.CODE_SQL, "查询删除文件失败：", err.Error())
	}
	for _, item := range upfiles {
		var err error
		switch item.Storage {
		case STORAGE_LOCAL:
			this.NewStorageLocal(localConfig).clean(item.Src)
		case STORAGE_QINIU:
			this.NewStorageQiniu(qiniuConfig).clean(item.Src)
		default:
		}
		if err != nil {
			try.Throwf(frame.CODE_SQL, "存储资源%d删除失败：", item.Id, err.Error())
		}
		this.OrmModel(item)
		item.Delete()
	}
}

type FindArgs struct {
	UserGroup string
	UserId    int64
	State     int
	BeginAt   string
	EndAt     string
	OrderCol  string
	OrderType string
}

// 上传图片数据查询
func (this *Component) Find(args *FindArgs, page, limit int, listPtr interface{}) int64 {
	db := this.OrmTable(this.NewFile())
	if args.UserId > 0 {
		db.And("`UserId`=? ", args.UserId)
	}
	if args.UserGroup != "" {
		db.And("`UserGroup`=?", args.UserGroup)
	}
	if args.State != 0 {
		db.And("State=?", args.State)
	}
	if args.BeginAt != "" {
		db.And("`CreateTime`>?", args.BeginAt)
	}
	if args.EndAt != "" {
		db.And("`CreateTime`<?", args.EndAt)
	}
	if args.OrderCol == "time" {
		if args.OrderType == "asc" {
			db.Asc("CreateTime")
		} else {
			db.Desc("CreateTime")
		}
	} else if args.OrderCol == "size" {
		if args.OrderType == "asc" {
			db.Asc("Size")
		} else {
			db.Desc("Size")
		}
	} else {
		db.Desc("Id")
	}
	return this.FindPage(db, listPtr, page, limit)
}
