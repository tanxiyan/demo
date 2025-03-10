package upload

import (
	"errors"
	"io"
	"os"
	"path"
	"time"

	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
	"gitee.com/go-mao/mao/libs/utils"
)

type storageLocal struct {
	config LocalConfig
}

func (this *Component) NewStorageLocal(config LocalConfig) *storageLocal {
	return &storageLocal{config: config}
}

func (this *storageLocal) mode() string {
	return STORAGE_LOCAL
}

// 上传到本地
func (this *storageLocal) upload(filename string, data io.Reader, size int64) (filesrc string, fileurl string) {
	ext := path.Ext(filename)
	filename = utils.Md5(filename, time.Now().UnixNano()) + ext
	subDir := time.Now().Format("20060102")
	dir := this.config.Dir + "/" + subDir
	if err := os.MkdirAll(dir, 0600); err != nil {
		try.Throwf(frame.CODE_WARN, "无法创建上传文件目录%s,错误：%s", dir, err.Error())
	}
	filesrc = dir + "/" + filename //本地文件路径
	f, e := os.Create(filesrc)
	if e != nil {
		try.Throwf(frame.CODE_WARN, "文件创建失败", e.Error())
		return
	}
	defer f.Close()
	if _, e := io.Copy(f, data); e != nil {
		try.Throwf(frame.CODE_WARN, "文件保存失败", e.Error())
	}
	fileurl = this.config.Domain + "/" + subDir + "/" + filename //网络访问路径
	return
}

// 清理
func (this *storageLocal) clean(filesrc string) error {
	if err := os.Remove(filesrc); err != nil {
		return errors.New("本地图片清理失败:" + err.Error())
	}
	return nil
}
