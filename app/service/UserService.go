package service

import (
	"log"
	"okc/app/model"
	"okc/utils"
)

// CheckExtensionCode
// 检查推广码是否被使用，如果被使用就重新生成，直到未使用
func CheckExtensionCode(extensionCode string) string {
	id, err := new(model.Users).QueryUserIdByExtensionCode(extensionCode)
	if err != nil {
		log.Println("CheckExtensionCode [ERROR] : ", err)
		return ""
	}
	if id != 0 {
		extensionCode = utils.GenerateRandExtensionCode(len(extensionCode))
		CheckExtensionCode(extensionCode)
	}
	return extensionCode
}
