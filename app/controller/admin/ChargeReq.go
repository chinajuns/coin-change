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

// ChargeReq
// 充值
type ChargeReq struct {
}

// QueryChargeReqList
// 获取充值列表
func (r *ChargeReq) QueryChargeReqList(c *gin.Context) {
	lang := c.GetHeader("lang")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	types := c.Query("type")
	account := c.Query("account_number")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	pageParseInt, _ := strconv.ParseInt(page, 10, 64)
	limitParseInt, _ := strconv.ParseInt(limit, 10, 64)
	data, total, err := new(model.ChargeReq).QueryChargeReqListPageByTimeAndTypeAndAccount(startTime, endTime, types, account, int(pageParseInt), int(limitParseInt))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.ChargeReq).QueryChargeReqListPageByTimeAndTypeAndAccount [ERROR] : %s", err))
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

// ChargePass
// 充值申请同意
func (r *ChargeReq) ChargePass(c *gin.Context) {
	lang := c.GetHeader("lang")
	id := c.PostForm("id")

	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

	chargeReq, err := new(model.ChargeReq).QueryChargeReqById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err.Error(),
		})
		return
	}

	if chargeReq.Id == 0 || chargeReq.Status != 1 {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": "没找到充值申请单或状态错误",
		})
		return
	}

	userWallet, err := new(model.UsersWallet).QueryUserWalletByUserIdAndCurrencyId(chargeReq.Uid, chargeReq.CurrencyId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.UsersWallet).QueryMapUserWalletByUserIdAndCurrencyId [ERROR] : %s", err))
		return
	}

	if userWallet.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": "用户钱包不存在",
		})
		return
	}

	res := service.ChangeUserWalletBalance(userWallet, 2, chargeReq.Amount, 200, "链上充币", false, 0, 0, "", 0, 0)
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
				"message": "充值失败",
			})
			return
		}
		break
	}

	chargeReq.Status = 2
	chargeReq.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	err = chargeReq.UpdateChargeReqById()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": "充值成功",
	})
	return
}

// ChargeReqRefuse
// 充值申请拒绝
func (r *ChargeReq) ChargeReqRefuse(c *gin.Context) {
	lang := c.GetHeader("lang")
	id := c.PostForm("id")
	marks := c.PostForm("marks")

	if id == "" || marks == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

	chargeReq, err := new(model.ChargeReq).QueryChargeReqById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.ChargeReq).QueryChargeReqById [ERROR] : %s", err))
		return
	}

	if chargeReq.Id == 0 || chargeReq.Status != 1 {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": "没找到充值申请单或状态错误",
		})
		return
	}

	chargeReq.Status = 3
	chargeReq.Remark = marks
	chargeReq.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	err = chargeReq.UpdateChargeReqById()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("chargeReq.UpdateChargeReqById [ERROR] : %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    400,
		"message": "拒绝成功",
	})
	return
}
