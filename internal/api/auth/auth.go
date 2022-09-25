package auth

import (
	"github.com/gin-gonic/gin"
	"godp/internal/api/types"
	"godp/internal/db"
	"godp/pkg/api/response"
	"godp/pkg/errorcode"
	"godp/pkg/permission"
	"gorm.io/gorm"
)

func UseAuthRouter(r *gin.RouterGroup) {
	r.POST("/auth/login", Login)
	r.POST("/auth/register", Register)
}

func Login(c *gin.Context) {
	param := &types.LoginParam{}
	err := c.BindJSON(&param)
	if err != nil {
		response.Error(c, errorcode.ErrInvalidParam)
	}
	username := param.Username
	password := param.Password
	if len(username) == 0 {
		response.Error(c, errorcode.ErrInvalidParam)
	}
	if len(password) == 0 {
		response.Error(c, errorcode.ErrInvalidParam)
	}

	err, user := db.SysUserDb.GetByUsername(username)
	if err == nil {
		if permission.PasswordVerify(password, user.PasswordHash) {
			tokenString, err := permission.GenerateToken(username, password)
			if err != nil {
				response.ErrorCustom(c, errorcode.ErrInvalidParam.Code(), "生成token失败")
			}
			response.Success(c, gin.H{
				"token": tokenString,
			})
		} else {
			response.ErrorCustom(c, errorcode.ErrInvalidParam.Code(), "密码错误")
		}
	} else if err == gorm.ErrRecordNotFound {
		response.ErrorCustom(c, errorcode.ErrInvalidParam.Code(), "用户不存在")
	} else {
		response.Error(c, errorcode.ErrDatabase)
	}
}

func Register(c *gin.Context) {
	param := &types.RegisterParam{}
	err := c.BindJSON(&param)
	if err != nil {
		response.Error(c, errorcode.ErrInvalidParam)
	}
	username := param.Username
	password := param.Password
	nickname := param.Nickname
	if len(username) == 0 {
		response.Error(c, errorcode.ErrInvalidParam)
	}
	if len(password) == 0 {
		response.Error(c, errorcode.ErrInvalidParam)
	}

	err, _ = db.SysUserDb.GetByUsername(username)
	if err == nil {
		response.ErrorCustom(c, errorcode.ErrInvalidParam.Code(), "username existed.")
	} else if err == gorm.ErrRecordNotFound {
		passwordHash, err := permission.PasswordHash(password)
		if err != nil {
			response.Error(c, errorcode.ErrEncrypt)
		}
		db.SysUserDb.Create(username, passwordHash, nickname)
	} else {
		response.Error(c, errorcode.ErrDatabase)
	}
}
