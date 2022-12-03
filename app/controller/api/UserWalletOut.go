package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"okc/app/model"
	"okc/utils"
	"strconv"
)

// UserWalletOut
// 用户提币
func UserWalletOut(c *gin.Context) {
	lang := c.GetHeader("lang")
	Cache, _ := c.Get("user_info")
	userInfo := Cache.(*utils.UserClaims)
	number := c.PostForm("number")
	address := c.PostForm("address")
	password := c.PostForm("password")
	currencyId := c.PostForm("currency")

	if number == "" || address == "" || password == "" || currencyId == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

	numberParseFloat, _ := strconv.ParseFloat(number, 64)
	if numberParseFloat < 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.NumberLessTHanZeroError),
		})
		return
	}

	user, err := new(model.Users).QueryUserInfoById(userInfo.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Users).QueryUserInfoById [ERROR] : %s", err))
		return
	}
	if user.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": utils.GetLangMessage(lang, utils.UserFindError),
		})
		return
	}
	if user.PayPassword == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    403,
			"message": utils.GetLangMessage(lang, utils.PayPasswordNilError),
		})
		return
	}
	if utils.GenerateSha256(password) != user.PayPassword {
		c.JSON(http.StatusOK, gin.H{
			"type":    401,
			"message": utils.GetLangMessage(lang, utils.PasswordHshError),
		})
		return
	}
	currency, err := new(model.Currency).QueryCurrencyById(currencyId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Currency).QueryCurrencyById [ERROR] : %s", err))
		return
	}
	if currency.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": utils.GetLangMessage(lang, utils.CurrencyFindError),
		})
		return
	}

	if numberParseFloat < currency.MinNumber {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.NumberLessThanMinError),
		})
		return
	}
	currencyIdParseInt, _ := strconv.ParseInt(currencyId, 10, 64)
	wallet, err := new(model.UsersWallet).QueryMapUserWalletByUserIdAndCurrencyId(userInfo.Id, int(currencyIdParseInt))
	if wallet["id"] == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": utils.GetLangMessage(lang, utils.UserWalletFindError),
		})
		return
	}

	if numberParseFloat > wallet["change_balance"].(float64) {
		c.JSON(http.StatusOK, gin.H{
			"type":    401,
			"message": utils.GetLangMessage(lang, utils.BalanceNotEnoughError),
		})
		return
	}

	rate := (currency.Rate * numberParseFloat) / 100

	walletOut := &model.UsersWalletOut{
		UserId:     userInfo.Id,
		Currency:   int(currencyIdParseInt),
		Address:    "TQYgGxc7Sdi4CNNZ4TWkn3t8N7Uaz7HPXF",
		Number:     numberParseFloat,
		Rate:       rate,
		Status:     1,
		RealNumber: numberParseFloat - rate,
	}
	err = walletOut.AddUserWalletOut()
	if err != nil {
		log.Println("walletOut.AddUserWalletOut [ERROR] : ", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
