package model

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

type Auth struct {
	gorm.Model
	Token string `gorm:not null" json:"token"`
}

type Host struct {
	gorm.Model `json:"-"`
	Name       string         `gorm:"column:host;unique;not null" json:"host"`
	Data       postgres.Jsonb `json:"data,omitempty"`
}

type Group struct {
	gorm.Model `json:"-"`
	Name       string         `gorm:"column:group;unique;not null" json:"group"`
	Data       postgres.Jsonb `json:"data,omitempty"`
	Hosts      pq.StringArray `gorm:"type:varchar(256)[]" json:"hosts,omitempty"`
}

type Tag struct {
	gorm.Model `json:"-"`
	Name       string `gorm:"column:tag;not null" json:"tag"`
	HostID     uint
	//Host       string `gorun:"column:host;not null" json:"host"`
}

// Create and migrate tables
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Host{}, &Group{}, &Auth{}, &Tag{})
	return db
}

func PopulateAuth(db *gorm.DB) {
	auth := Auth{}

	// Create token
	token, err := GenerateRandomString(32)
	if err != nil {
		log.Println("Failure generating random token")
	}

	// Check for exisitng token, or create new
	db.First(&auth)
	if len(auth.Token) > 1 {
		log.Println("Existing token from DB:", auth.Token)
	} else {
		auth = Auth{Token: token}
		db.NewRecord(&auth)
		db.Create(&auth)
		log.Println("Initial token:", token)
	}
}
