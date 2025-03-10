package upload

import "io"

const (
	STORAGE_LOCAL = "local" //上传到本地
	STORAGE_QINIU = "qiniu" //上传到七牛
)

const (
	FILE_STATE_USED   = 1  //已经使用
	FILE_STATE_NOTUSE = -1 //未使用
)

type UserInterface interface {
	GetUserId() int64
	GetUserGroup() string
}

type FileStorage interface {
	upload(filename string, data io.Reader, size int64) (src string, url string)
	clean(src string) error
	mode() string
}
