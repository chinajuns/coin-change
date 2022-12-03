package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
	"strconv"
)

// Users
// 用户管理
type Users struct {
}

// QueryUserList
// 获取用户列表
func (u *Users) QueryUserList(c *gin.Context) {
	lang := c.GetHeader("lang")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("page", "10")
	account := c.Query("account")
	name := c.Query("name")
	risk := c.DefaultQuery("risk", "-2")

	pageParseInt, _ := strconv.ParseInt(page, 10, 64)
	limitParseInt, _ := strconv.ParseInt(limit, 10, 64)
	riskParseInt, _ := strconv.ParseInt(risk, 10, 64)
	data, total, err := new(model.Users).QueryUsersListPageByAccountAndNameAndRisk(account, name, int(riskParseInt), int(pageParseInt), int(limitParseInt))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Users).QueryUsersListPageByAccountAndNameAndRisk [ERROR] : %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    data,
		"total":   total,
		"message": utils.GetLangMessage(lang, utils.Success),
	})

	return
}

// QueryUserWalletList
// 获取用户钱包列表
func (u *Users) QueryUserWalletList(c *gin.Context) {
	lang := c.GetHeader("lang")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	userId := c.Query("user_id")

	if userId == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

	pageParseInt, _ := strconv.ParseInt(page, 10, 64)
	limitParseInt, _ := strconv.ParseInt(limit, 10, 64)

	data, total, err := new(model.UsersWallet).QueryUserWalletListPageByUserId(userId, int(pageParseInt), int(limitParseInt))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Users).QueryUsersListPageByAccountAndNameAndRisk [ERROR] : %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    data,
		"total":   total,
		"message": utils.GetLangMessage(lang, utils.Success),
	})

	return
}
