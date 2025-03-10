package main

import (
	"embed"
	"waimai_api/globals"
	"waimai_api/modules/base"
	"waimai_api/modules/base/service"
	"waimai_api/modules/cms"
	"waimai_api/modules/coupon"
	"waimai_api/modules/points"
	"waimai_api/modules/printer"
	"waimai_api/modules/shop"
	"waimai_api/modules/sys"

	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
)

//go:embed views/*
var UiEmbed embed.FS

func main() {
	globals.UiEmbed = UiEmbed
	sev := frame.NewServer(new(App))
	sev.Run()
}

type App struct {
}

// 程序启动时执行
func (this *App) Init(m *frame.WebEngine) {

}

// 配置文件
func (this *App) ConfigFile() string {
	return "./config.ini"
}

// 模块
func (this *App) Modules() []frame.ModuleInterface {
	return []frame.ModuleInterface{
		base.New(),
		cms.New(),
		coupon.New(),
		shop.New(),
		sys.New(),
		points.New(),
		printer.New(),
	}
}

// 安装执行脚本
func (this *App) Install() error {
	return nil
}

// 权限验证
func (this *App) CheckPermission(w *frame.Webline, code, name string) {
	adminEntity := service.GetLoginAdmin(w)
	if adminEntity.GetUserId() == 10000 {
		return
	}
	permissionService := service.NewPermissionService(w.Taskline)
	if !permissionService.Check(adminEntity, code) {
		try.Throw(globals.CODE_NOT_AUTH, "没有权限")
	}
}

// 版本信息
// test2
func (this *App) Version() string {
	return globals.Version()
}
