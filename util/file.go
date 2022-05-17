package util

import (
	"io/ioutil"
	"log"
)

//使用ioutil.WriteFile方式写入文件,是将[]byte内容写入文件,如果content字符串中没有换行符的话，默认就不会有换行符
func WriteWithIoutil(name, content string) error {
	data := []byte(content)
	if err := ioutil.WriteFile(name, data, 0644); err == nil {
		log.Println("写入文件成功:", name)
	} else {
		return err
	}
	return nil
}
