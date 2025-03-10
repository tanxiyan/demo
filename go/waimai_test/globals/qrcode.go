package globals

import (
	"fmt"
	"gitee.com/go-mao/mao/libs/utils"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

var QrcodePath = "/qrcode/:key"

var qrcodesMap = sync.Map{}

// 输出二维码
func WriteQrcodeHandle(c *gin.Context) {
	key := c.Param("key")
	val, ok := qrcodesMap.Load(key)
	if !ok {
		c.String(200, "二维码已过期")
		return
	}
	data, err := qrcode.Encode(val.(string), qrcode.Medium, 256)
	if err != nil {
		c.String(200, "二维码生成失败")
		return
	}
	c.Data(200, "image/jpeg", data)
	c.Abort()
}

// 创建二维码链接
func MakeQrcodeLink(data string) string {
	key := utils.Sha1(data, time.Now().UnixNano())
	qrcodesMap.Store(key, data)
	return fmt.Sprintf("/qrcode/%s", key)
}
