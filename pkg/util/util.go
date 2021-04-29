package util

import (
	"log"
	"operator-test-framework/pkg/api"
	"os"
	"strings"
)

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsDir(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		log.Println(err)
		return false
	}
	return s.IsDir()
}

func ConvertStrToPara(input string, para []api.Parameter) []api.Parameter {
	input = input + ","
	array := strings.Split(input, ",")
	for _, str := range array {
		if !strings.Contains(str, "=") {
			continue
		}
		array1 := strings.Split(str, "=")
		if len(array) != 2 {
			continue
		}
		exist := false
		for i, _ := range para {
			if para[i].Name == array1[0] {
				para[i].Value = array1[1]
				exist = true
			}
		}
		if exist == false {
			para = append(para, api.Parameter{
				Name:  array1[0],
				Value: array1[1],
			})
		}
	}
	return para
}
