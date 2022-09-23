package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"godp/internal/api/types"
	"godp/internal/db"
	_pojo "godp/internal/pojo"
	"godp/internal/process"
	"godp/pkg/errorcode"
	"godp/pkg/helper"
	"gorm.io/gorm"
	"strconv"

	"godp/pkg/page"
	"godp/pkg/response"
)

func UseProjectRouter(r *gin.RouterGroup) {
	r.GET("/project/page", ProjectPage)
	r.GET("/project/:id", ProjectGet)
	r.POST("/project", ProjectPost)
	r.PUT("/project/:id", ProjectPut)
	r.PUT("/project/:id/config", ProjectUpdateConfig)
	r.DELETE("/project/:id", ProjectDelete)
}

func ProjectPage(c *gin.Context) {
	param := &types.ProjectQueryParam{}
	err := c.ShouldBindQuery(&param)
	if helper.HandleError(c, err) {
		return
	}
	err, objList, total := db.ProjectInfoDb.Page(*param)
	if helper.HandleError(c, err) {
		return
	}
	pageResult := page.PageResultWrapper(param.CurPage, param.Limit, objList, total)
	response.Success(c, pageResult)
	return
}

func ProjectGet(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Param("id"))
	if helper.HandleError(c, err) {
		return
	}
	err, projectInfo := db.ProjectInfoDb.GetById(uint(projectId))
	if helper.HandleError(c, err) {
		return
	}
	response.Success(c, projectInfo)
}

func ProjectPost(c *gin.Context) {
	saveOrUpdate(c, 0)
}

func ProjectPut(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Param("id"))
	if helper.HandleError(c, err) {
		return
	}
	saveOrUpdate(c, uint(projectId))
}

func ProjectUpdateConfig(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Param("id"))
	if helper.HandleError(c, err) {
		return
	}
	param := make(map[string]interface{})
	c.BindJSON(&param)
	config := param["config"].(string)
	projectConfig := _pojo.ProjectConfig{}
	err = json.Unmarshal([]byte(config), &projectConfig)
	if helper.HandleError(c, err) {
		return
	}
	// 获取内网地址
	for i, ipAddress := range projectConfig.IPAddressArr {
		err, intranetIp := process.GetIntranetIp(projectConfig, ipAddress.IP)
		if err != nil {
			continue
		}
		fmt.Println(ipAddress.IP + ":" + intranetIp)
		projectConfig.IPAddressArr[i].IntranetIP = intranetIp
		if i == 0 {
			if projectConfig.MysqlIp == "" || projectConfig.MysqlIp == "127.0.0.1" {
				projectConfig.MysqlIp = intranetIp
			}
			if projectConfig.RedisIp == "" || projectConfig.RedisIp == "127.0.0.1" {
				projectConfig.RedisIp = intranetIp
			}
		}
	}
	db.ProjectInfoDb.UpdateConfig(uint(projectId), projectConfig)
}

func ProjectDelete(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Param("id"))
	if helper.HandleError(c, err) {
		return
	}
	err = db.ProjectInfoDb.DeleteById(uint(projectId))
	if helper.HandleError(c, err) {
		return
	}
	response.Success(c, "")
}

func saveOrUpdate(c *gin.Context, projectId uint) {
	param := &db.ProjectInfo{}
	err := c.BindJSON(&param)
	if helper.HandleError(c, err) {
		return
	}
	param.ID = projectId
	err, _ = db.ProjectInfoDb.GetByProjectApp(param.ProjectApp, param.ID)
	if err == nil {
		response.ErrorCustom(c, errorcode.ErrInvalidParam.Code(), "项目编码已经存在")
		return
	} else if err == gorm.ErrRecordNotFound {
		err = db.ProjectInfoDb.CreateOrUpdate(*param)
		if helper.HandleError(c, err) {
			return
		}
		response.Success(c, "")
	} else {
		response.Error(c, errorcode.ErrDatabase)
	}
}
