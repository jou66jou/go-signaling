package db

import (
	"errors"
	"fmt"
	"log"
	"p2p/internal/config"
	"strings"
	"time"

	"github.com/cenk/backoff"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// InitDatabase 從 config 初始化
func InitDatabase(cfg *config.Configuration, name string) (db *gorm.DB, err error) {
	for _, database := range cfg.Databases {
		if strings.EqualFold(database.Name, name) {
			db, err = setupDatabase(database)
			if err != nil {
				return nil, err
			}
			return db, nil
		}
	}
	return nil, errors.New("configs no match db name : " + name)
}

func setupDatabase(database config.Database) (*gorm.DB, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&multiStatements=true", database.Username, database.Password, database.Address, database.DBName)
	var db *gorm.DB
	var err error
	err = backoff.Retry(func() error {
		db, err = gorm.Open("mysql", connectionString)
		if err != nil {
			log.Printf("main: mysql open failed: %v \n", err)
			return err
		}
		err = db.DB().Ping()
		if err != nil {
			log.Printf("main: mysql ping error: %v \n", err)
			return err
		}
		return nil
	}, bo)

	if err != nil {
		return nil, err
	}

	db.LogMode(true)

	db.DB().SetMaxIdleConns(150)
	db.DB().SetMaxOpenConns(300)
	db.DB().SetConnMaxLifetime(14400 * time.Second)

	return db, nil
}
