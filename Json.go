package main

import (
	"strconv"
)

type JsonValueType interface{}

// Json 结构体定义
type Json struct {
	keys   []string
	values []JsonValueType
}

// Append 方法：向Json对象添加键值对
func (js *Json) AppendString(key string, value string) {
	js.keys = append(js.keys, key)
	var value_ JsonValueType
	value_ = value
	js.values = append(js.values, value_)
}
func (js *Json) AppendJson(key string, js_ Json) {
	js.keys = append(js.keys, key)
	var value_ JsonValueType
	value_ = js_
	js.values = append(js.values, value_)
}
func (js *Json) AppendInt(key string, val int) {
	js.keys = append(js.keys, key)
	var value_ JsonValueType
	value_ = val
	js.values = append(js.values, value_)
}
func (js *Json) AppendBool(key string, val bool) {
	js.keys = append(js.keys, key)
	var value_ JsonValueType
	value_ = val
	js.values = append(js.values, value_)
}
// Get 方法：将Json对象转换为JSON格式字符串
func (js *Json) Get() string {
	res := "{"
	length := len(js.keys)
	for i := 0; i < length; i++ {
		if i != 0 {
			res += ","
		}
		value := js.values[i]
		valueStr := ""
		//
		if str, ok := value.(string); ok {
			valueStr = "\"" + str + "\""
		} else if js_, ok := value.(Json); ok {
			valueStr = js_.Get()
		} else if val, ok := value.(int); ok {
			valueStr = strconv.Itoa(val)
		} else if val, ok := value.(bool); ok {
			if val {
				valueStr = "true"
			} else {
				valueStr = "false"
			}
		}
		//
		res += "\"" + js.keys[i] + "\":" + valueStr
	}
	res += "}"
	return res
}
