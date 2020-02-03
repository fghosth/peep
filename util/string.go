package util

import "errors"

//首字母转小写
func FUPer(str string) (string, error) {
	errEmpty := errors.New("字符串为空")
	v := []byte(str)
	if len(v) == 0 {
		return "", errEmpty
	}
	if v[0] < 97 {
		v[0] += 32
	}
	return string(v), nil
}
