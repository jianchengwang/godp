package assetsHost

import (
	"github.com/gin-gonic/gin"
	"godp/internal/api/types"
	"godp/internal/db"
	"godp/pkg/api/helper"
	"godp/pkg/api/page"
	"godp/pkg/api/response"
	"godp/pkg/errorcode"
	"gorm.io/gorm"
	"strconv"
)

func UseAssetsHostRouter(r *gin.RouterGroup) {
	r.GET("/assetsHost/page", HostPage)
	r.GET("/assetsHost/:id", HostGet)
	r.POST("/assetsHost", HostPost)
	r.PUT("/assetsHost/:id", HostPut)
	r.DELETE("/assetsHost/:id", HostDelete)
}

func HostPage(c *gin.Context) {
	param := &types.AssetsHostQueryParam{}
	err := c.ShouldBindQuery(&param)
	if helper.HandleError(c, err) {
		return
	}
	err, objList, total := db.AssetsHostDb.Page(*param)
	if helper.HandleError(c, err) {
		return
	}
	pageResult := page.PageResultWrapper(param.CurPage, param.Limit, objList, total)
	response.Success(c, pageResult)
}

func HostGet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if helper.HandleError(c, err) {
		return
	}
	err, projectInfo := db.AssetsHostDb.GetById(uint(id))
	if helper.HandleError(c, err) {
		return
	}
	response.Success(c, projectInfo)
}

func HostPost(c *gin.Context) {
	saveOrUpdate(c, 0)
}

func HostPut(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if helper.HandleError(c, err) {
		return
	}
	saveOrUpdate(c, uint(id))
}

func HostDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if helper.HandleError(c, err) {
		return
	}
	err = db.AssetsHostDb.DeleteById(uint(id))
	if helper.HandleError(c, err) {
		return
	}
	response.Success(c, "")
}

func saveOrUpdate(c *gin.Context, id uint) {
	param := &db.AssetsHost{}
	err := c.BindJSON(&param)
	if helper.HandleError(c, err) {
		return
	}
	param.ID = id
	err, _ = db.AssetsHostDb.GetByIP(param.IP, param.ID)
	if err == nil {
		response.ErrorCustom(c, errorcode.ErrInvalidParam.Code(), "ip exited.")
		return
	} else if err == gorm.ErrRecordNotFound {
		err = db.AssetsHostDb.SaveOrUpdate(*param)
		if helper.HandleError(c, err) {
			return
		}
		response.Success(c, "")
	} else {
		response.Error(c, errorcode.ErrDatabase)
	}
}
