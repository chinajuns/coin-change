package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/app/service"
	"okc/utils"
)

// LeverClose
// 合约平仓
func LeverClose(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)
	id := c.PostForm("id")

	lever, err := new(model.LeverTransaction).QueryLeverTransactionById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.LeverTransaction).QueryLeverTransactionById [ERROR] : %s", err))
		return
	}
	// 判断是否为空
	if lever.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": "数据未找到",
		})
		return
	}
	// 判断是否同一个用户
	if lever.UserId != userInfo.Id {
		c.JSON(http.StatusOK, gin.H{
			"type":    401,
			"message": "无权操作",
		})
		return
	}
	// 判断状态
	if lever.Status != 1 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "交易状态异常,请勿重复提交",
		})
		return
	}

	// 封装
	err = service.CheckLeverClose(lever)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
