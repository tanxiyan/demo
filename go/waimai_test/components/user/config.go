package user

type Config struct {
	MaxErrLogin  int          //最大错误次数
	AllowDevices []DeviceInfo //登录终端配置
}

type DeviceInfo struct {
	Name       string `name:"终端名称"`
	Label      string `name:"终端标识"`
	MaxSession int    `name:"会话时常/分"`
}

// 验证终端类型
func (this *Config) GetDevice(device string) *DeviceInfo {
	if device == "" {
		return &DeviceInfo{}
	}
	for _, item := range this.AllowDevices {
		if item.Label == device {
			return &item
		}
	}
	return nil
}
