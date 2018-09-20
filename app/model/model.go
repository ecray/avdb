package model

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

type Auth struct {
	gorm.Model
	Token string `gorm:not null" json:"token"`
}

type Host struct {
	gorm.Model
	Name string         `gorm:"unique;not null" json:"host"`
	Data postgres.Jsonb `json:"data"`
}

type Group struct {
	gorm.Model
	Name  string         `gorm:"unique;not null" json:"group"`
	Data  postgres.Jsonb `json:"data"`
	Hosts pq.StringArray `gorm:"type:varchar(256)[]" json:"hosts"`
}

// Create and migrate tables
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Host{}, &Group{}, &Auth{})
	return db
}

func PopulateAuth(db *gorm.DB) {
	auth := Auth{Token: "XxXtestXxX"}
	db.NewRecord(&auth)
	db.Create(&auth)
}
