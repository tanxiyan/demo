package recharge_withdraw

import (
	"gitee.com/go-mao/mao/frame"
)

type Component struct {
	*frame.Taskline
	config  *Config
	prefix  string
	handler RechargeWithdrawhHandlerInteface
}

// prefix=表前缀
func New(t *frame.Taskline, handler RechargeWithdrawhHandlerInteface) *Component {
	object := new(Component)
	object.Taskline = t
	object.handler = handler
	object.prefix = handler.Prefix()
	return object
}

func (this *Component) GetConfig() *Config {
	return this.handler.GetConfig()
}

func (this *Component) Models() []frame.ModelInterface {
	models := make([]frame.ModelInterface, 0)
	models = append(models, this.NewRecharge())
	models = append(models, this.NewWithdraw())
	return models
}
