package coin

import (
	"encoding/json"
	"fmt"
	"testing"

	"gitee.com/go-mao/mao/libs/binding"
)

func aTestCoin(t *testing.T) {
	var coin Coin
	var str = "90"
	err := json.Unmarshal([]byte(str), &coin)
	if err != nil {
		t.Log("整型反序列化解析错误:", err.Error())
	} else {
		t.Log("整型反序列化成功:", coin)
	}
	coinType = COIN_TYPE_TWO
	str = "90.11"
	err = json.Unmarshal([]byte(str), &coin)
	if err != nil {
		t.Log("浮点型反序列化解析错误:", err.Error())
	} else {
		t.Log("浮点型反序列化成功:", coin)
	}
	fmt.Printf("哈哈：%d", coin)

}

func bTestBind(t *testing.T) {
	var req struct {
		Number Coin
	}
	data := map[string]interface{}{
		"Number": 98,
	}
	coinType = COIN_TYPE_INT
	err := binding.Bind(&req, data).Error()
	if err != nil {
		t.Log("整型数据绑定失败", err.Error())
	} else {
		t.Log("整型数据绑定成功", req.Number)
	}

	//
	var req2 struct {
		Number Coin `valid:"gt=100:错误"`
	}
	var ms interface{} = &req2.Number
	if _, ok := ms.(json.Unmarshaler); ok {
		t.Log("Coin已实现json.Unmarshaler")
	}
	data = map[string]interface{}{
		"Number": "98.33333",
	}
	coinType = COIN_TYPE_TWO
	err = binding.Bind(&req2, data).Error()
	if err != nil {
		t.Log("浮点型数据绑定失败", err.Error())
	} else {
		t.Log("浮点型数据绑定成功", int64(req2.Number))
	}
	fmt.Printf("哈哈2：%d", req2.Number)
}

func TestCharge(t *testing.T) {
	coinType = COIN_TYPE_FOUR
	c := Coin(10)
	fmt.Println(int64(c))
	fmt.Println(CoinByCharge(c, 0.56))

}
