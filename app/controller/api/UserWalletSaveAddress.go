package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"okc/app/model"
	"okc/utils"
)

// UserWalletSaveAddress
// 用户钱包地址保存
func UserWalletSaveAddress(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)
	address := c.PostForm("address")
	if address == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    401,
			"message": utils.GetLangMessage(lang, utils.AddressInputError),
		})
		return
	}
	userAddressModel := new(model.UsersAddress)
	usersAddress, err := userAddressModel.QueryFindByUserId(userInfo.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersAddress).QueryFindByUserIdAndAddress [ERROR] : %s", err))
		return
	}
	log.Println("usersAddress:", usersAddress)
	// 检查是否已经有地址
	if usersAddress.Id == 0 {
		// 新增
		err = userAddressModel.AddUserAddressByUserId(userInfo.Id, address)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    500,
				"message": utils.GetLangMessage(lang, utils.ServerError),
			})
			_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersAddress).QueryFindByUserIdAndAddress [ERROR] : %s", err))
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"type":    "ok",
			"message": utils.GetLangMessage(lang, utils.Success),
		})
		return
	} else {
		// 更新
		err = userAddressModel.UpdateUserAddressByUserId(userInfo.Id, address)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    500,
				"message": utils.GetLangMessage(lang, utils.ServerError),
			})
			_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersAddress).UpdateUserAddressByUserId [ERROR] : %s", err))
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"type":    "ok",
			"message": utils.GetLangMessage(lang, utils.Success),
		})
		return
	}

}
