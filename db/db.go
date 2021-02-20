package db

import (
	"github.com/fatih/color"
	"github.com/inblack67/GQLGenAPI/mymodels"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// PgConn ...
	PgConn *gorm.DB
)

// ConnectDB ...
func ConnectDB() (*gorm.DB){
	dsn := "host=localhost user=postgres password=postgres dbname=gographql port=5432"
	var err error
	PgConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		panic(err)
	}
	PgConn.AutoMigrate(&mymodels.User{}, &mymodels.Story{})
	color.Green("Postgres is here")
	return PgConn
}