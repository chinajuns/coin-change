package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"okc/app/model"
	"okc/app/service"
	"okc/utils"
	"strconv"
	"time"
)

// UserEmailRegister
// 用户邮箱注册
func UserEmailRegister(c *gin.Context) {
	// 语言
	lang := c.PostForm("lang")
	// 邮箱
	email := c.PostForm("user_string")
	// 验证码
	code := c.PostForm("code")
	// 密码
	password := c.PostForm("password")
	// 确认密码
	rePassword := c.PostForm("re_password")
	// 推荐码
	extensionCode := c.PostForm("extension_code")
	// 注册区号
	areaCodeId := c.DefaultPostForm("area_code_id", "0")
	// 注册区号
	areaCode := c.DefaultPostForm("area_code", "0")
	// 用户模型
	user := new(model.Users)

	defer func() {
		err := recover()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    "500",
				"message": utils.GetLangMessage(lang, utils.ServerError),
			})
			log.Println(err)
			_ = utils.WriteErrorLog(fmt.Sprintf("UserEmailRegister [ERROR] : %s", err))
			return
		}
	}()

	// 验证必填参数
	if email == "" || code == "" || password == "" ||
		rePassword == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    "401",
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

	// 检查密码一致
	if password != rePassword {
		c.JSON(http.StatusOK, gin.H{
			"type":    "401",
			"message": utils.GetLangMessage(lang, utils.PasswordHshError),
		})
		return
	}

	// 检查密码长度
	if len(password) < 6 || len(password) > 16 {
		c.JSON(http.StatusOK, gin.H{
			"type":    "401",
			"message": utils.GetLangMessage(lang, utils.PasswordLenError),
		})
		return
	}

	redis := utils.RedisConnect()
	defer func() {
		err := redis.Close()
		if err != nil {
			_ = utils.WriteErrorLog(err)
		}
	}()
	// 检查验证码是否匹配或存在
	reply, err := redis.Do("GET", "EMAIL_CODE_"+email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    "401",
			"message": utils.GetLangMessage(lang, utils.CodeError),
		})
		_ = utils.WriteErrorLog(err)
		return
	}
	if reply == nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    "401",
			"message": utils.GetLangMessage(lang, utils.CodeError),
		})
		return
	}
	if code != string(reply.([]byte)) {
		c.JSON(http.StatusOK, gin.H{
			"type":    "401",
			"message": utils.GetLangMessage(lang, utils.CodeError),
		})
		return
	}
	// 检查账号是否被注册
	isAdd, err := new(model.Users).CheckUserExistsByPhoneAndEmail(email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    "500",
			"message": utils.GetLangMessage(lang, utils.ServerError) + err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("user.CheckUserExistsByPhoneAndEmail [ERROR] : %s", err))
		return
	}
	if isAdd {
		c.JSON(http.StatusOK, gin.H{
			"type":    "401",
			"message": utils.GetLangMessage(lang, utils.AccountExistsError),
		})
		return
	}
	// 检查推广码
	if extensionCode != "" {
		parentId, err := new(model.Users).QueryUserIdByExtensionCode(extensionCode)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    "500",
				"message": utils.GetLangMessage(lang, utils.ServerError) + err.Error(),
			})
			_ = utils.WriteErrorLog(fmt.Sprintf("user.QueryUserIdByExtensionCode [ERROR] : %s", err))
			return
		}
		if parentId == 0 {
			c.JSON(http.StatusOK, gin.H{
				"type":    "401",
				"message": utils.GetLangMessage(lang, utils.PcodeError),
			})
			return
		}
		user.ParentId = parentId
	}

	// 注册
	user.AccountNumber = email
	user.Email = email
	user.Password = utils.GenerateSha256(password)
	user.PayPassword = utils.GenerateSha256(password)
	user.WalletPwd = utils.GenerateSha256(password)
	user.Time = int(time.Now().Unix())
	user.ExtensionCode = service.CheckExtensionCode(utils.GenerateRandExtensionCode(4))

	idInt, _ := strconv.Atoi(areaCodeId)
	codeInt, _ := strconv.Atoi(areaCode)
	user.AreaCodeId = idInt
	user.AreaCode = codeInt

	err = user.AddUserByEmail()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    "500",
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("user.AddUserByEmail [ERROR] : %s", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": utils.GetLangMessage(lang, utils.Success),
		"data":    gin.H{},
	})
	return
}
