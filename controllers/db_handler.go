package controllers

import (
	"database/sql"
	"log"

	"github.com/go-redis/redis/v8"
)

var client *redis.Client

func SetupRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // no password set
		DB:       0,                // default DB
	})
}

func SetupDatabase() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db_tugasexplorasi3?parseTime=true&loc=Asia%2FJakarta")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
