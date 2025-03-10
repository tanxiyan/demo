package upload

import (
	"io"
	"path"
	"strings"

	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
	"gitee.com/go-mao/mao/libs/utils"
)

// 上传处理域
type uploader struct {
	filename string
	data     io.Reader
	size     int64
	*Component
}

func (this *Component) newUploader(data io.Reader, filename string, size int64) *uploader {
	up := &uploader{
		Component: this,
		data:      data,
		size:      size,
		filename:  filename,
	}
	return up
}

// 检查文件类型
func (this *uploader) checkfile() {
	config := this.Config().(*Config)
	ext := strings.ToLower(path.Ext(this.filename))
	if !strings.Contains(config.AllowExts, ext) {
		try.Throw(frame.CODE_WARN, "禁止上传"+ext+"文件")
	}
	//检查文件大小
	if this.size/1024 > config.MaxSize*1024 {
		try.Throwf(frame.CODE_WARN, "不允许上传超过%dM的文件", config.MaxSize)
	}
}

// 根据配置，自动保存
func (this *uploader) SaveAuto(operator UserInterface) string {
	config := this.Config().(*Config)
	switch config.Storage {
	case STORAGE_LOCAL:
		return this.SaveStorage(operator, this.NewStorageLocal(config.LocalConfig))
	case STORAGE_QINIU:
		return this.SaveStorage(operator, this.NewStorageQiniu(config.QiniuConfig))
	default:
		try.Throw(frame.CODE_WARN, "上传失败", "不支持的上传类型")
	}
	return ""
}

// 根据上传引擎上传
func (this *uploader) SaveStorage(operator UserInterface, storage FileStorage) string {
	this.checkfile()
	src, url := storage.upload(this.filename, this.data, this.size)
	fileEntity := this.NewFile()
	fileEntity.UserId = operator.GetUserId()
	fileEntity.UserGroup = operator.GetUserGroup()
	fileEntity.Src = src
	fileEntity.Url = url
	fileEntity.Storage = storage.mode()
	fileEntity.Md5 = utils.Md5(url)
	fileEntity.Size = this.size
	fileEntity.State = FILE_STATE_NOTUSE
	if !fileEntity.Create() {
		try.Throwf(frame.CODE_WARN, "文件%s上传失败，无法保存到数据库", this.filename)
	}
	return url
}
