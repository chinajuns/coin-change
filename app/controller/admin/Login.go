package admin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"okc/app/model"
	"okc/utils"

	"github.com/gin-gonic/gin"
)

// Login
// 登录
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	// verfy_code := c.PostForm("code")
	log.Println("username : ", username)
	admin, err := new(model.Admin).QueryAdminByUsername(username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage("服务器错误", utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Admin).QueryAdminByUsername [ERROR] : %s", err))
		return
	}
	if admin.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": utils.GetLangMessage("账号不存在", utils.UserFindError),
		})
		return
	}
	//fmt.Printf("utils.GenerateMd5(\"123456\"): %v\n", utils.GenerateMd5("123456"))
	if utils.GenerateHashPassword(password, true) != admin.Password {
		c.JSON(http.StatusOK, gin.H{
			"type":    401,
			"message": utils.GetLangMessage("密码错误", utils.PasswordError),
		})
		return
	}

	token, err := utils.GenerateToken(admin.Id, admin.Username, admin.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage("token错误", utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("utils.GenerateToken [ERROR] : %s", err))
		return
	}

	adminRole, err := new(model.AdminRole).QueryAdminRoleById(admin.RoleId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage("", utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.AdminRole).QueryAdminRoleById [ERROR] : %s", err))
		return
	}
	right := adminRole.Right
	log.Println("right :", right)
	var rightPares []interface{}
	_ = json.Unmarshal([]byte(right), &rightPares)

	c.JSON(http.StatusOK, gin.H{
		"type": "ok",
		"data": gin.H{
			"token": token,
			"user":  admin,
			"role":  rightPares,
		},
		"message": utils.GetLangMessage("success", utils.Success),
	})
	return
}
