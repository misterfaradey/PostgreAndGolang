package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type StorageConf struct {
	DbHost     string
	DbPort     string
	DbName     string
	DbScheme   string
	DbUser     string
	DbPassword string
}

type DB interface {
	Connect(config *StorageConf) error
	Close() error
	HealthChecker(logger *log.Logger)
	InitDB() error

	Wallet
}

type db struct {
	db *sql.DB
}

func NewCarStorage() DB {
	return &db{}
}

func (c *db) Connect(config *StorageConf) error {

	s := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable connect_timeout=10",
		config.DbHost,
		config.DbPort,
		config.DbName,
		config.DbUser,
		config.DbPassword,
	)

	db, err := sql.Open("postgres", s)
	if err != nil {
		return err
	}

	c.db = db

	//todo add stmt

	db.SetMaxOpenConns(500)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(30 * time.Second)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	return db.PingContext(ctx)
}

func (c *db) Close() error {
	return c.db.Close()
}

func (c *db) HealthChecker(logger *log.Logger) {

	for {
		time.Sleep(time.Minute)

		ctx, _ := context.WithTimeout(context.Background(), time.Second*15)
		err := c.db.PingContext(ctx)
		if err != nil {
			logger.Println("db Ping", err)
		}
	}
}

//это временное решение
func (c *db) InitDB() error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

	rows, err := c.db.QueryContext(ctx, `SELECT 1 FROM "test".test LIMIT 1;`)
	if err == nil {
		rows.Close()
		return nil
	}

	_, err = c.db.ExecContext(ctx, initDB)
	return err
}
