package upload

import (
	"context"
	"errors"
	"io"
	"path"
	"time"

	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
	"gitee.com/go-mao/mao/libs/utils"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type storageQiniu struct {
	config QiniuConfig
	dir    string
}

func (this *Component) NewStorageQiniu(config QiniuConfig) *storageQiniu {
	return &storageQiniu{config: config}
}

func (this *Component) StorageQiniuOnDir(config QiniuConfig, dir string) *storageQiniu {
	return &storageQiniu{config: config, dir: dir}
}

func (this *storageQiniu) mode() string {
	return STORAGE_QINIU
}

func (this *storageQiniu) qiniuConfig() (*qbox.Mac, storage.Config) {
	conf := this.config
	mac := qbox.NewMac(conf.AccessKey, conf.Secret)
	var zone *storage.Zone
	switch conf.Zone {
	case "hn": //华南
		zone = &storage.ZoneHuanan
	case "hb": //华北
		zone = &storage.ZoneHuabei
	case "hd": //华东
		zone = &storage.ZoneHuadong
	case "as0": //东南亚
		zone = &storage.Zone_as0
	case "cn-east-2": //华东-浙江
		zone = &storage.ZoneFogCnEast1
	case "na0": //北美
		zone = &storage.Zone_na0
	default:
		try.Throw(frame.CODE_WARN, "无效存储区域")
	}
	cfg := storage.Config{
		Zone:          zone,  // 空间对应的机房
		UseHTTPS:      false, // 是否使用https域名
		UseCdnDomains: true,  // 上传是否使用CDN上传加速
	}
	return mac, cfg
}

// 上传到七牛
func (this *storageQiniu) upload(filename string, data io.Reader, size int64) (filesrc string, fileurl string) {
	ext := path.Ext(filename)
	filename = utils.Md5(filename, time.Now().UnixNano()) + ext
	conf := this.config
	putPolicy := storage.PutPolicy{
		Scope: conf.Bucket,
	}
	mac, cfg := this.qiniuConfig()
	upToken := putPolicy.UploadToken(mac)
	formUploader := storage.NewFormUploader(&cfg) // 构建表单上传的对象
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{},
	}
	if this.dir != "" {
		filename = this.dir + "/" + filename
	}
	filesrc = filename
	err := formUploader.Put(context.Background(), &ret, upToken, filename, data, size, &putExtra)
	if err != nil {
		try.Throw(frame.CODE_WARN, "上传失败", err.Error())
	}
	fileurl = conf.Domain + "/" + filename //网络访问路径
	return
}

// 清理
func (this *storageQiniu) clean(filesrc string) error {
	conf := this.config
	deleteOps := []string{storage.URIDelete(conf.Bucket, filesrc)}
	mac, cfg := this.qiniuConfig()
	bucketManager := storage.NewBucketManager(mac, &cfg)
	_, err := bucketManager.Batch(deleteOps)
	if err != nil {
		return errors.New("七牛图片清理失败" + err.Error())
	}
	return nil
}
