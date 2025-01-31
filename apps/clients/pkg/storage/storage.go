package storage

import (
	"fmt"
	"log"

	"github.com/owjoel/client-factpack/apps/clients/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type storage struct {
	Client ClientInterface
}

var (
	db *storage
)

func Init() {
	_db, err := gorm.Open(mysql.Open(GetDSN()), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	log.Printf("Connected to DB")
	db = &storage{Client: &ClientStorage{_db}}
}

// username:password@protocol(address)/dbname?param=value
func GetDSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)
	return dsn
}

func GetInstance() *storage {
	return db
}
