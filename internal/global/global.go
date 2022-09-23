package global

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var JwtToken = "Authorization"
var JwtSecret = "godp@1234"
var JwtIssuer = "godp"

var CI_IP = ""
var CI_IntranetIP = ""
var CI_PORT = 22
var CI_USER = ""
var CI_PASSWORD = ""
var CI_DEPLOYDIR = ""

var RepoAccessKeyId = ""
var RepoAccessKeySecret = ""
var RepoRegionId = ""
var RepoNamespace = ""

var DB *gorm.DB
var Route *gin.Engine
