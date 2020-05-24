package models

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"set-flags/pkg/setting"
	"time"
)

var db *gorm.DB

// Base contains common columns for all tables.
type Base struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time  `json:"created_at"`
	UpdateAt  time.Time  `json:"update_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func InitDB() {
	var (
		err                                        error
		dbType, dbName, user, password, host, port string
	)

	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}
	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	port = sec.Key("PORT").String()

	db, err = gorm.Open(dbType, fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		host,
		port,
		user,
		dbName,
		password))

	if err != nil {
		log.Println(err)
	}

	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// Migrate the schema
	db.AutoMigrate(&Flag{}, &Asset{}, &Evidence{})
}

func CloseDB() {
	defer db.Close()
}
