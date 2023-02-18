package database

import (
	"davisbento/golang-encoder/domain"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type Database struct {
	DB            *gorm.DB
	Dsn           string
	DsnTest       string
	DBType        string
	DBTypeTest    string
	Debug         bool
	AutoMigration bool
	Env           string
}

func NewDB() *Database {
	return &Database{}
}

func NewDBTest() *gorm.DB {
	dbInstance := NewDB()
	dbInstance.Env = "test"
	dbInstance.DBTypeTest = "sqlite3"
	dbInstance.DsnTest = ":memory:"
	dbInstance.AutoMigration = true
	dbInstance.Debug = true

	conn, err := dbInstance.Connect()

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	return conn
}

func (db *Database) Connect() (*gorm.DB, error) {
	var err error

	if db.Env == "test" {
		db.DB, err = gorm.Open(db.DBTypeTest, db.DsnTest)
	} else {
		db.DB, err = gorm.Open(db.DBType, db.Dsn)
	}

	if err != nil {
		return nil, err
	}

	if db.Debug {
		db.DB.LogMode(true)
	}

	if db.AutoMigration {
		db.DB.AutoMigrate(&domain.Video{}, &domain.Job{})
	}

	return db.DB, nil
}
