package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/app/service"
	"okc/utils"
	"strconv"
)

// UserWalletChange
// 用户钱包转划
func UserWalletChange(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)

	types := map[string]int{
		"legal":  1,
		"lever":  3,
		"micro":  4,
		"change": 2,
	}
	momo := map[string]string{
		"legal":  "c2c",
		"lever":  "合约",
		"micro":  "期货",
		"change": "币币",
	}
	fromArr := map[string]int{
		"legal":  9,
		"lever":  16,
		"micro":  16,
		"change": 12,
	}
	toArr := map[string]int{
		"legal":  10,
		"lever":  15,
		"micro":  15,
		"change": 11,
	}
	currencyId := c.PostForm("currency_id")
	number := c.PostForm("number")
	fromField := c.PostForm("from_field")
	toField := c.PostForm("to_field")

	if currencyId == "" || number == "" || fromField == "" || toField == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

	fromAccountLogType := fromArr[fromField] // 转出
	toAccountLogType := toArr[toField]       // 转入

	// 备注
	remark := fmt.Sprintf("%s 转划 %s", momo[fromField], momo[toField])

	if fromField == "lever" {
		count, err := new(model.LeverTransaction).QueryLeverTransactionCountByUserIdAndInStatus(userInfo.Id, "3, 4")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    500,
				"message": utils.GetLangMessage(lang, utils.ServerError),
			})
			_ = utils.WriteErrorLog(fmt.Sprintf("new(model.LeverTransaction).QueryLeverTransactionCountByUserIdAndInStatus [ERROR] %s \n", err))
			return
		}
		if count > 0 {
			c.JSON(http.StatusOK, gin.H{
				"type":    400,
				"message": "您有正在进行中的杆杠交易,不能进行此操作",
			})
			return
		}
	}
	// 用户钱包
	userWallet, err := new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId(userInfo.Id, currencyId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId [ERROR] %s \n", err))
		return
	}
	if userWallet.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": utils.GetLangMessage(lang, utils.UserWalletFindError),
		})
		return
	}
	numberParseFloat, _ := strconv.ParseFloat(number, 64)

	res := service.ChangeUserWalletBalance(userWallet, types[fromField], numberParseFloat*-1, fromAccountLogType, remark, false, 0, 0, "", 0, 0)
	switch res.(type) {
	case error:
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": res.(error),
		})
		return
		break
	case string:
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": res.(string),
		})
		return
		break
	case bool:
		if res.(bool) == false {
			c.JSON(http.StatusOK, gin.H{
				"type":    400,
				"message": "转划失败",
			})
			return
		}
		break
	}
	userWallet.Refresh()
	res = service.ChangeUserWalletBalance(userWallet, types[toField], numberParseFloat, toAccountLogType, remark, false, 0, 0, "", 0, 0)
	switch res.(type) {
	case error:
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": res.(error),
		})
		return
		break
	case string:
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": res.(string),
		})
		return
		break
	case bool:
		if res.(bool) == false {
			c.JSON(http.StatusOK, gin.H{
				"type":    400,
				"message": "转划失败",
			})
			return
		}
		break
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return

}
