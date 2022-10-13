package api

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"godp/internal/api/assetsHost"
	"godp/internal/api/auth"
	"godp/internal/api/batchScript"
	"godp/internal/api/middleware"
	ws "godp/internal/api/ws"
)

func InitRoute() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.MaxMultipartMemory = 32 << 20 // 32 MiB

	v1 := r.Group("/api/v1/")
	{
		auth.UseAuthRouter(v1)
	}
	v1.Use(middleware.JWT())

	assetsHost.UseAssetsHostRouter(v1)
	batchScript.UseBatchScriptRouter(v1)
	ws.UseWsSshRouter(v1)

	// docs
	r.Use(static.Serve("/docs", static.LocalFile("./docs/", false)))
	// frontend
	r.Use(static.Serve("/", static.LocalFile("./frontend/dist/", false)))

	r.Run(":8081")
	return r
}
