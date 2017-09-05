package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type MainData struct {
	Command string
	User    string
	Port    string
}

type MainBody struct {
	Ip   string
	Data []MainData
}

func main() {

	//增加标志位
	var targetBit bool = true
	var s, inits MainBody

	str := `{"Ip":"10.1.1.199","data":[{"Command":"redis","User":"root","Port":"8086"},{"Command":"nginx","User":"root","Port":"80"},{"Command":"master","User":"root","Port":"2502"}]}`
	initStr := `{"Ip":"10.1.1.88", "data": [{"Command":"nginx","User":"root","Port": "8202"}]}`

	if err := json.Unmarshal([]byte(str), &s); err != nil {
		fmt.Println("json err:", err)
	}
	if err := json.Unmarshal([]byte(initStr), &inits); err != nil {
		fmt.Println("json err:", err)
	}

	for _, comeinData := range s.Data {
		//如果数据相同，就不用跳出循环，并将标志位设置为false.当不一样时，标志位的结果是true, 这样循环后会被打印
		for _, initData := range inits.Data {
			if reflect.DeepEqual(comeinData, initData) {
				targetBit = false
				break
			}
		}
		//如果标志位为true,刚打开不一致的数据
		if targetBit {
			fmt.Println(comeinData)
		}
		//设置标志位为默认值
		targetBit = true
	}
}
