package api

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

// CoinTradeCancel
// 币币订单取消
func CoinTradeCancel(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)
	id := c.Query("id")
	coin, err := new(model.CoinTrade).QueryCoinTradeById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.CoinTrade).QueryCoinTradeById [ERROR] : %s", err))
		return
	}
	if coin.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "订单未找到",
		})
		return
	}
	if coin.UId != userInfo.Id {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "请求异常",
		})
		return
	}
	if coin.Status != 1 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "状态异常",
		})
		return
	}
	switch coin.Type {
	case 1:
		TargetPriceParseStr := strconv.FormatFloat(coin.TargetPrice, 'f', 4, 64)
		TradeAmountParseStr := strconv.FormatFloat(coin.TradeAmount, 'f', 4, 64)
		err := service.UnUserBuyCoint(coin.UId, coin.CurrencyId, coin.LegalId, TargetPriceParseStr, TradeAmountParseStr)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    400,
				"message": err,
			})
			return
		}
		// 更新订单
		coin.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		coin.Status = 3
		err = coin.UpdateCoinTradeById()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    500,
				"message": utils.GetLangMessage(lang, utils.ServerError),
			})
			_ = utils.WriteErrorLog(fmt.Sprintf("coin.UpdateCoinTradeById [ERROR] : %s", err))
			return
		}
		break
	case 2:
		TradeAmountParseStr := strconv.FormatFloat(coin.TradeAmount, 'f', 4, 64)
		err := service.UnUserSellCoin(coin.UId, coin.CurrencyId, TradeAmountParseStr)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    400,
				"message": err,
			})
			return
		}
		// 更新订单
		coin.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		coin.Status = 3
		err = coin.UpdateCoinTradeById()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    500,
				"message": utils.GetLangMessage(lang, utils.ServerError),
			})
			_ = utils.WriteErrorLog(fmt.Sprintf("coin.UpdateCoinTradeById [ERROR] : %s", err))
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
