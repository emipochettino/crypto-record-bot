package persistance

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Db *gorm.DB
}

func NewDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("crypto_record.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	autoMigrate(db)

	return db
}

//autoMigrate Migrate the schema
func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&AlertDAO{},
	)
	if err != nil {
		panic("failed to migrate schema")
	}
}
