package admin

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/app/model"
	"okc/utils"
	"strconv"
	"strings"
)

// AdminRole
// 管理员角色
type AdminRole struct {
}

// AddAdminRole
// 添加管理员角色
func (a *AdminRole) AddAdminRole(c *gin.Context) {
	name := c.PostForm("name")
	right := c.PostForm("right")
	isSuper := c.DefaultPostForm("isSuper", "0")

	if name == "" || right == "" || isSuper == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage("", utils.ParameterError),
		})
		return
	}

	adminRole, err := new(model.AdminRole).QueryAdminRoleByName(name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.AdminRole).QueryAdminRoleByName [ERROR] : %s", err))
		return
	}

	if adminRole.Id > 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "角色名称已存在",
		})
		return
	}
	rightSlice := strings.Split(right, ",")
	rightSliceJson, _ := json.Marshal(rightSlice)
	err = new(model.AdminRole).AddAdminRole(name, string(rightSliceJson), isSuper)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": "添加失败: " + err.Error(),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.AdminRole).AddAdminRole [ERROR] : %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"message": "添加成功",
	})
	return

}

// EditAdminRole
// 编辑管理员角色
func (a *AdminRole) EditAdminRole(c *gin.Context) {
	name := c.PostForm("name")
	right := c.PostForm("right")
	isSuper := c.DefaultPostForm("isSuper", "0")

	if name == "" || right == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "参数不能为空",
		})
		return
	}

	// 查询管理员角色是否存在
	adminRole, err := new(model.AdminRole).QueryAdminRoleByName(name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage("", utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.AdminRole).QueryAdminByUsername [ERROR] : %s ", err))
		return
	}
	if adminRole.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "管理员角色不存在",
		})
		return
	}

	// 更新管理员角色
	if strings.Index(right, ",") != -1 {
		rightSlice := strings.Split(right, ",")
		rightSliceJson, _ := json.Marshal(rightSlice)
		adminRole.Right = string(rightSliceJson)
	} else {
		rightSliceJson, _ := json.Marshal([]string{right})
		adminRole.Right = string(rightSliceJson)
	}

	isSuperParseInt, _ := strconv.ParseInt(isSuper, 10, 64)
	adminRole.Name = name
	adminRole.IsSuper = int(isSuperParseInt)
	err = adminRole.UpdateAdminRoleById()
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

// DeleteAdminRole
// 删除管理员角色
func (a *AdminRole) DeleteAdminRole(c *gin.Context) {
	name := c.PostForm("name")

	if name == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "参数不能为空",
		})
		return
	}

	// 查询管理员角色是否存在
	adminRole, err := new(model.AdminRole).QueryAdminRoleByName(name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage("", utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.AdminRole).QueryAdminByUsername [ERROR] : %s ", err))
		return
	}
	if adminRole.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": "管理员角色不存在",
		})
		return
	}

	// 删除管理员
	err = adminRole.DeleteAdminRoleById()
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

// QueryAdminRoleList
// 获取管理员角色列表
func (a *AdminRole) QueryAdminRoleList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageParseInt, _ := strconv.ParseInt(page, 10, 64)
	limitParseInt, _ := strconv.ParseInt(limit, 10, 64)

	data, total, err := new(model.AdminRole).QueryAdminRoleListPage(int(pageParseInt), int(limitParseInt))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err,
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.AdminRole).QueryAdminListPage [ERROR] : %s ", err))
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

// QueryAdminRoleAll
// 获取管理员所有角色
func (a *AdminRole) QueryAdminRoleAll(c *gin.Context) {

	data, err := new(model.AdminRole).QueryAdminRoleAll()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": err,
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("new(model.AdminRole).QueryAdminRoleAll [ERROR] : %s ", err))
		return
	}

	slice := make([]map[string]interface{}, 0)
	for _, v := range data {
		slice = append(slice, map[string]interface{}{
			"value": v.Id,
			"label": v.Name,
		})

	}
	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    slice,
		"message": "获取成功",
	})

	return
}
