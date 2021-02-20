package db

import (
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/inblack67/GQLGenAPI/mymodels"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// PgConn ...
	PgConn *gorm.DB
)

// ConnectDB ...
func ConnectDB() (*gorm.DB){

	dsn := "host=localhost user=postgres password=postgres dbname=gographql port=5432"

	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), 
		logger.Config{
			SlowThreshold: time.Microsecond, 	// to make it all queries log  
			LogLevel:      logger.Info, 
			Colorful:      true,         
		},
	)

	PgConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil{
		panic(err)
	}

	PgConn.AutoMigrate(&mymodels.User{}, &mymodels.Story{})
	color.Green("Postgres is here")
	return PgConn
}