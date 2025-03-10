package coin

import "github.com/shopspring/decimal"

//根据输入的base金额的charge比例数量，向上取整
func CoinByCharge(coinNumber Coin, charge float64) int64 {

	return decimal.NewFromInt(int64(coinNumber)).Mul(decimal.NewFromFloat(charge)).Add(decimal.NewFromFloat(0.5)).IntPart()

	//return int64(float64(coinNumber)*charge + 0.5)
}
