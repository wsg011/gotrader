package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// 获取精度数据
func DecimalMath(str string) int32 {
	s := strings.Split(str, ".")
	if len(s) == 2 {
		f := strings.Split(s[1], "")
		for i := len(f); i > 0; i-- { //循环去掉末尾的0
			if f[i-1] != strconv.Itoa(0) {
				return int32(i)
			}
		}
	} else {
		return 0
	}
	return 0
}

func ByteToMap(msg []byte) (map[string]interface{}, error) {
	if res, err := ByteToInterface(msg); err != nil {
		return nil, err
	} else {
		switch msgType := res.(type) {
		case map[string]interface{}:
			return res.(map[string]interface{}), nil
		default:
			return nil, fmt.Errorf("%v", msgType)
		}
	}
}

func ByteToInterface(byteMsg []byte) (interface{}, error) {
	var tempMap interface{}
	d := json.NewDecoder(bytes.NewReader(byteMsg))
	d.UseNumber()
	err := d.Decode(&tempMap)
	if err != nil {
		return nil, err
	}
	return tempMap, nil
}
