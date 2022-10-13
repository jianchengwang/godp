package batchScript

import (
	"github.com/gin-gonic/gin"
	"godp/internal/api/types"
	"godp/internal/db"
	"godp/pkg/api/helper"
	"godp/pkg/api/page"
	"godp/pkg/api/response"
	"strconv"
)

func UseBatchScriptRouter(r *gin.RouterGroup) {
	r.GET("/batchScript/page", ScriptPage)
	r.GET("/batchScript/:id", ScriptGet)
	r.POST("/batchScript", ScriptPost)
	r.PUT("/batchScript/:id", ScriptPut)
	r.DELETE("/batchScript/:id", ScriptDelete)
}

func ScriptPage(c *gin.Context) {
	param := &types.BatchScriptsQueryParam{}
	err := c.ShouldBindQuery(&param)
	if helper.HandleError(c, err) {
		return
	}
	err, objList, total := db.BatchScriptDb.Page(*param)
	if helper.HandleError(c, err) {
		return
	}
	pageResult := page.PageResultWrapper(param.CurPage, param.Limit, objList, total)
	response.Success(c, pageResult)
}

func ScriptGet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if helper.HandleError(c, err) {
		return
	}
	err, projectInfo := db.BatchScriptDb.GetById(uint(id))
	if helper.HandleError(c, err) {
		return
	}
	response.Success(c, projectInfo)
}

func ScriptPost(c *gin.Context) {
	saveOrUpdate(c, 0)
}

func ScriptPut(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if helper.HandleError(c, err) {
		return
	}
	saveOrUpdate(c, uint(id))
}

func ScriptDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if helper.HandleError(c, err) {
		return
	}
	err = db.BatchScriptDb.DeleteById(uint(id))
	if helper.HandleError(c, err) {
		return
	}
	response.Success(c, "")
}

func saveOrUpdate(c *gin.Context, id uint) {
	param := &db.BatchScript{}
	err := c.BindJSON(&param)
	if helper.HandleError(c, err) {
		return
	}
	param.ID = id
	err = db.BatchScriptDb.SaveOrUpdate(*param)
	if helper.HandleError(c, err) {
		return
	}
	response.Success(c, "")
}
