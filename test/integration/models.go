package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name      string
	Gender    string
	Languages []Language `gorm:"many2many:user_languages;"`
	Note      string

	Profile Profile
}

type Profile struct {
	gorm.Model
	UserId  uint64
	Address string
}

type Language struct {
	gorm.Model
	Name string
}

type Product struct {
	gorm.Model
	Name        string
	Description string
}

var (
	DB      gorm.DB
	devMode bool
	dbname  string
	dbuser  string
	dbpwd   string
)

func PrepareDB() {
	var err error

	// Be able to start a server for develop test
	dbuser, dbpwd = "qor", "qor"
	if devMode {
		dbname = "qor_integration"
	} else {
		dbname = "qor_integration_test"
	}

	DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", dbuser, dbpwd, dbname))
	if err != nil {
		panic(err)
	}

	setupDb(!devMode) // Don't drop table in dev mode
}

func getTables() []interface{} {
	return []interface{}{
		&User{},
		&Product{},
		&Profile{},
		&Language{},
	}
}

func setupDb(dropBeforeCreate bool) {
	tables := getTables()

	for _, table := range tables {
		if dropBeforeCreate {
			if err := DB.DropTableIfExists(table).Error; err != nil {
				panic(err)
			}
		}

		if err := DB.AutoMigrate(table).Error; err != nil {
			panic(err)
		}
	}
}
