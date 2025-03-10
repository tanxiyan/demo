package coin

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
)

//货币类型
type CoinType int8

const (
	COIN_TYPE_INT  CoinType = 1 //整形货币
	COIN_TYPE_TWO  CoinType = 2 //两位小数货币
	COIN_TYPE_FOUR CoinType = 4 //4位小数
	COIN_TYPE_SIX  CoinType = 6 //6位小数
)

//平台货币类型
var coinType CoinType

//设置货币类型
func SetCoinType(t CoinType) {
	if coinType != 0 {
		panic("不可重复设置货币类型")
	}
	coinType = t
}

//获取货币类型
func GetCoinType() CoinType {
	return coinType
}

//获取设置的小数比例
func GetRate() int64 {
	switch coinType {
	case COIN_TYPE_TWO: //输出两位小数
		return 100
	case COIN_TYPE_FOUR:
		return 10000
	case COIN_TYPE_SIX:
		return 1000000
	default: //输出整形
		return 1
	}
}

//根据货币类型返回相应的数除比例
func GetRateByCoinType(t CoinType) int64 {
	switch t {
	case COIN_TYPE_TWO: //输出两位小数
		return 100
	case COIN_TYPE_FOUR:
		return 10000
	case COIN_TYPE_SIX:
		return 1000000
	default: //输出整形
		return 1
	}
}

//根据默认配置获取货币,浮点型转整形
func GetCoin(d float64) Coin {
	switch coinType {
	case COIN_TYPE_TWO: //输出两位小数
		return Coin(d * 100)
	case COIN_TYPE_FOUR:
		return Coin(d * 10000)
	case COIN_TYPE_SIX:
		return Coin(d * 1000000)
	default: //输出整形
		return Coin(d)
	}
}

//平台货币
type Coin int64

func (self Coin) String() string {
	switch coinType {
	case COIN_TYPE_TWO: //输出两位小数
		return fmt.Sprintf(`%.2f`, float64(self)/100)
	case COIN_TYPE_FOUR:
		return fmt.Sprintf(`%.4f`, float64(self)/10000)
	case COIN_TYPE_SIX:
		return fmt.Sprintf(`%.6f`, float64(self)/1000000)
	default: //输出整形
		return fmt.Sprintf(`%d`, self)
	}
}

//反序列化
func (self *Coin) UnmarshalJSON(data []byte) error {
	data = bytes.Replace(data, []byte(`"`), []byte(""), -1) //避免出现字符串类型数据
	switch coinType {
	case COIN_TYPE_TWO: //输出两位小数
		number, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return err
		}
		*self = Coin(decimal.NewFromFloat(number).Mul(decimal.NewFromFloat(100)).IntPart())
	case COIN_TYPE_FOUR: //输出4位小数
		number, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return err
		}
		*self = Coin(decimal.NewFromFloat(number).Mul(decimal.NewFromFloat(10000)).IntPart())
	case COIN_TYPE_SIX: //输出6位小数
		number, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return err
		}
		*self = Coin(decimal.NewFromFloat(number).Mul(decimal.NewFromFloat(1000000)).IntPart())
	default: //输出整形
		number, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		*self = Coin(number)
	}
	return nil
}

//
func (self Coin) MarshalJSON() ([]byte, error) {
	return []byte(self.String()), nil
}

//数据格式验证
func (self Coin) Validate() error {
	return nil
}

func (self Coin) Int64() int64 {
	return int64(self)
}
