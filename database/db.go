package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

/*
	1. creating a global DB variable that can be accessed throughout the app
	2. Init function for DB
*/

type DB struct {
	*sql.DB
}

var Db *DB
var Rdb *redis.Client

func NewDBconn() *DB {
	// open the DB driver by its name and DB sourse
	db, err := sql.Open("postgres", os.Getenv("PG_DB"))
	if err != nil {
		log.Panicf("Error while opening up the postgres driver %+v", err)
	}
	return &DB{db}
}
