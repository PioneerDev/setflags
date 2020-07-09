package models

import (
	"fmt"
	"log"
	"set-flags/pkg/setting"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"

	// postgres dirver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

// Base contains common columns for all tables.
type Base struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time  `json:"created_at"`
	UpdateAt  time.Time  `json:"update_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// InitDB init db
func InitDB() {
	var (
		err                                  error
		dbType, dbName, user, password, host string
		port                                 int
	)

	dbType = setting.GetConfig().DataBase.Type
	dbName = setting.GetConfig().DataBase.Name
	user = setting.GetConfig().DataBase.User
	password = setting.GetConfig().DataBase.Password
	host = setting.GetConfig().DataBase.Host
	port = setting.GetConfig().DataBase.Port

	db, err = gorm.Open(dbType, fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
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
	db.AutoMigrate(&Flag{}, &Asset{}, &Evidence{}, &User{}, &Witness{}, &Payment{})
}

// CloseDB close db connection
func CloseDB() {
	defer db.Close()
}
