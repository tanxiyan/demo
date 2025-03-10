package globals

import "embed"

//ui资源
var UiEmbed embed.FS

const (
	CODE_NOT_LOGIN    = 401 //未登录
	CODE_NOT_AUTH     = 403 //没有权限
	CODE_API_NOT_AUTH = 802 //接口验证失败
)

var (
	FINAL_SERVICE_MARKET_API_ADDR = "https://api.zhiyuan2022.shop" //服务市场总图地址
	SERVICE_MARKET_API_ADDR       = ""                             //服务市场地址
)

const (
	RUN_MODE_SAAS = "saas"
)
