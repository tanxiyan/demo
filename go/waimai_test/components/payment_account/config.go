package payment_account

type Config struct {
	Group              string   //
	GroupName          string   //账号分组名称
	AccountTypes       []string //支持的支付方式
	AllowUpdateWaitDay int      //绑定账号可修改时间
}

func (this *Config) CheckAccountType(typ string) bool {
	for _, item := range this.AccountTypes {
		if item == typ {
			return true
		}
	}
	return false
}
