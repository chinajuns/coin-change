package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
)

// UserRechargeAddress
// 用户充值地址
func UserRechargeAddress(c *gin.Context) {
	lang := c.GetHeader("lang")
	currencyId := c.PostForm("currency")
	if currencyId == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

	//currency, err := new(model.Currency).QueryCurrencyById(currencyId)
	//if err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"type":    500,
	//		"message": utils.GetLangMessage(lang, utils.ServerError),
	//	})
	//	_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Currency).QueryCurrencyById [ERROR] : %s", err))
	//	return
	//}
	//if currency.Id == 0 {
	//	c.JSON(http.StatusOK, gin.H{
	//		"type":    404,
	//		"message": utils.GetLangMessage(lang, utils.Error),
	//	})
	//	return
	//}

	address := make(map[string]interface{})

	setting := new(model.Setting)
	address["OMNI"], _ = setting.QueryValueByKey("OMNI", "OMNI")
	address["ERC20"], _ = setting.QueryValueByKey("OMNI", "ERC20")
	address["TRC20"], _ = setting.QueryValueByKey("TRC20", "TRC20")

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    address,
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
