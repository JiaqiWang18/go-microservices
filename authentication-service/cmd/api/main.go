package main

import (
	"authentication/data"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")
	conn := connectToDB()
	if conn == nil {
		log.Panic("Could not connect to database")
		return
	}
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		db, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres is not ready yet. Waiting...")

			counts++
		} else {
			log.Println("Connected to database")
			return db
		}
		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Retrying in 2 seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}
