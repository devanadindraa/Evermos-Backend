package database

import (
	"fmt"
	"log"
	"time"

	"github.com/devanadindraa/Evermos-Backend/utils/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(conf *config.Config) (gormDb *gorm.DB, err error) {

	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Name,
	)

	var (
		gormDB *gorm.DB
	)

	for i := 0; i < 20; i++ {
		gormDB, err = getGormDB(connStr)
		if err == nil {
			break
		}

		log.Print(("Database not ready yet, retrying in 10 seconds..."))
		time.Sleep(10 * time.Second)
	}

	return gormDB, nil
}

func getGormDB(connStr string) (gormDB *gorm.DB, err error) {
	gormDB, err = gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}
