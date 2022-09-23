package main

import (
	"godp/internal/api"
	_ "godp/internal/config"
	"godp/internal/db"
	"godp/internal/global"
)

func main() {
	global.DB = db.InitDB()
	global.Route = api.InitRoute()
}
