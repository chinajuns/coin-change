package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/app/service"
	"okc/utils"
	"strconv"
	"time"
)

// UserWalletOut
// 用户提币申请
type UserWalletOut struct {
}

// QueryUserWalletOutList
// 获取用户提币列表
func (w *UserWalletOut) QueryUserWalletOutList(c *gin.Context) {
	lang := c.GetHeader("lang")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	account := c.DefaultQuery("account_number", "")

	pageParseInt, _ := strconv.ParseInt(page, 10, 64)
	limitParseInt, _ := strconv.ParseInt(limit, 10, 64)

	data, total, err := new(model.UsersWalletOut).QueryPagesByAccount(account, int(pageParseInt), int(limitParseInt))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersWalletOut).QueryPagesByAccount [ERROR] : %s", err))
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

// QueryUserWalletOutInfo
// 获取用户提币申请信息
func (w *UserWalletOut) QueryUserWalletOutInfo(c *gin.Context) {
	lang := c.GetHeader("lang")
	id := c.Query("id")

	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

	userWalletOut, err := new(model.UsersWalletOut).QueryUserWalletOutById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersWalletOut).QueryUserWalletOutById [ERROR] : %s", err))
		return
	}

	if userWalletOut.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": "没有找到此申请",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    userWalletOut,
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}

// QueryUserWalletOutPass
// 同意用户提币申请
func (w *UserWalletOut) QueryUserWalletOutPass(c *gin.Context) {
	lang := c.GetHeader("lang")
	id := c.PostForm("id")
	method := c.PostForm("method")
	notes := c.PostForm("notes")
	verificationcode := c.PostForm("verificationcode")

	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

	userWalletOut, err := new(model.UsersWalletOut).QueryUserWalletOutById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersWalletOut).QueryUserWalletOutById [ERROR] : %s", err))
		return
	}

	if userWalletOut.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": "未找到此申请单",
		})
		return
	}

	number := userWalletOut.Number
	userId := userWalletOut.UserId
	currency := userWalletOut.Currency

	userWallet, err := new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId(userId, currency)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("ew(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId [ERROR] : %s", err))
		return
	}

	if userWallet.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": "未找到此用户钱包",
		})
		return
	}

	if method == "done" {

		res := service.ChangeUserWalletBalance(userWallet, 2, number*-1, 100, "提币成功", true, 0, 0, "", 0, 0)
		switch res.(type) {
		case error:
			c.JSON(http.StatusOK, gin.H{
				"type":    400,
				"message": res.(error).Error(),
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
					"message": "操作失败",
				})
				return
			}
			break
		}

		userWalletOut.Status = 2
		userWalletOut.Notes = notes
		userWalletOut.Verificationcode = verificationcode
		userWalletOut.UpdateTime = int(time.Now().Unix())
		err := userWalletOut.UpdateUserWalletOut()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    400,
				"message": "操作失败 : " + err.Error(),
			})
			_ = utils.WriteErrorLog(fmt.Sprintf("userWalletOut.UpdateUserWalletOut [ERROR] : %s", err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"type":    "ok",
			"message": "操作成功",
		})
		return
	}

	// 锁定余额减少
	res := service.ChangeUserWalletBalance(userWallet, 2, number*-1, 101, "提币失败,锁定余额减少", true, 0, 0, "", 0, 0)
	switch res.(type) {
	case error:
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": res.(error).Error(),
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
				"message": "操作失败",
			})
			return
		}
		break
	}

	// 锁定余额撤回
	res = service.ChangeUserWalletBalance(userWallet, 2, number, 101, "提币失败,锁定余额撤回", false, 0, 0, "", 0, 0)
	switch res.(type) {
	case error:
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": res.(error).Error(),
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
				"message": "操作失败",
			})
			return
		}
		break
	}

	// 更新状态
	userWalletOut.Status = 3
	userWalletOut.Notes = notes
	userWalletOut.Verificationcode = verificationcode
	userWalletOut.UpdateTime = int(time.Now().Unix())
	err = userWalletOut.UpdateUserWalletOut()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "操作失败 : " + err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("userWalletOut.UpdateUserWalletOut [ERROR] : %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": "操作成功",
	})
	return
}
