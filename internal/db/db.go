package db

import (
	"fmt"
	"godp/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

func InitDB() *gorm.DB {
	mysqlConfig := config.Config.Mysql
	charset := "utf8mb4"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		mysqlConfig.Username,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.Database,
		charset)
	// open db
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		println(err)
	}
	// autoMigrate
	_ = db.AutoMigrate(
		&SysUser{},
		&SysRole{},
		&AssetsHost{},
	)
	return db
}
