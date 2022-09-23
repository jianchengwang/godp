package db

import (
	"godp/internal/global"
	"time"
)

type SysUser struct {
	BaseModel
	Username      string    `json:"username"`
	PasswordHash  string    `json:"passwordHash"`
	Nickname      string    `json:"nickname"`
	LastLoginTime time.Time `json:"lastLoginTime"`
}

type SysUserDbStruct struct{}

var SysUserDb = new(SysUserDbStruct)

func (SysUserDb *SysUserDbStruct) GetByUsername(username string) (error, SysUser) {
	var user SysUser
	err := global.DB.Where("username = ?", username).First(&user).Error
	return err, user
}

func (SysUserDb *SysUserDbStruct) Create(username string, passwordHash string, nickname string) {
	global.DB.Create(&SysUser{
		Username:      username,
		PasswordHash:  passwordHash,
		Nickname:      nickname,
		LastLoginTime: time.Now(),
	})
}
