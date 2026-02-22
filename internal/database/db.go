package database

import (
	"context"
	"os"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var conn *pgxpool.Pool

func Connection() error {
	var err error
    envData := os.Getenv("DATABASE_URL")
    
	log.Println("Connecting to:", envData)
	if envData == "" {
	    return errors.New("DATABASE_URL is not set")
	}
	conn, err = pgxpool.New(context.Background(), envData)
	if err != nil {
		return err
	}
	return conn.Ping(context.Background())
}

func CloseConnection() {

	if conn != nil {
		conn.Close()
	}

}

func GetDB() *pgxpool.Pool {
	return conn
}
