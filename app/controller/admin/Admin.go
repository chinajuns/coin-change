package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
	"strconv"
)

// Admin
// 管理员
type Admin struct {
}

// AddAdmin
// 添加管理员
func (a *Admin) AddAdmin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	roleId := c.PostForm("role_id")

	if username == "" || password == "" || roleId == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "参数不能为空",
		})
		return
	}

	// 查询管理员是否存在
	admin, err := new(model.Admin).QueryAdminByUsername(username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage("", utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Admin).QueryAdminByUsername [ERROR] : %s ", err))
		return
	}
	if admin.Id > 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "用户名已存在",
		})
		return
	}

	// 创建管理员
	roleIdParseInt, _ := strconv.ParseInt(roleId, 10, 64)
	admins := &model.Admin{
		Username: username,
		Password: utils.GenerateHashPassword(password, true),
		RoleId:   int(roleIdParseInt),
	}
	err = admins.AddAdmin()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "注册失败",
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("admins.AddAdmin() [ERROR] : %s ", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": "注册成功",
	})
	return
}

// EditAdmin
// 编辑管理员
func (a *Admin) EditAdmin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	roleId := c.PostForm("role_id")

	if username == "" || roleId == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "参数不能为空",
		})
		return
	}

	// 查询管理员是否存在
	admin, err := new(model.Admin).QueryAdminByUsername(username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage("", utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Admin).QueryAdminByUsername [ERROR] : %s ", err))
		return
	}
	if admin.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "用户不存在",
		})
		return
	}

	// 更新管理员
	roleIdParseInt, _ := strconv.ParseInt(roleId, 10, 64)
	admin.RoleId = int(roleIdParseInt)
	if password != "" {
		admin.Password = utils.GenerateHashPassword(password, true)
	}
	err = admin.UpdateAdminById()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": fmt.Sprintf("更新失败 : %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": "更新成功",
	})
	return
}

// DeleteAdmin
// 删除管理员
func (a *Admin) DeleteAdmin(c *gin.Context) {
	username := c.PostForm("username")
	cache, _ := c.Get("user_info")
	userClami := cache.(utils.UserClaims)
	if username == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "参数不能为空",
		})
		return
	}

	// 查询管理员是否存在
	admin, err := new(model.Admin).QueryAdminByUsername(username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage("", utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Admin).QueryAdminByUsername [ERROR] : %s ", err))
		return
	}
	if admin.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "用户不存在",
		})
		return
	}
	if admin.Username == userClami.Account {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "自己不能删除自己",
		})
		return
	}

	// 删除管理员
	err = admin.DeleteAdminById()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": fmt.Sprintf("删除失败 : %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": "删除成功",
	})
	return
}

// QueryAdminList
// 获取管理员列表
func (a *Admin) QueryAdminList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageParseInt, _ := strconv.ParseInt(page, 10, 64)
	limitParseInt, _ := strconv.ParseInt(limit, 10, 64)

	data, total, err := new(model.Admin).QueryAdminListPage(int(pageParseInt), int(limitParseInt))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err,
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.Admin).QueryAdminListPage [ERROR] : %s ", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    data,
		"total":   total,
		"message": "获取成功",
	})

	return
}
