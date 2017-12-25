package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/acoshift/configfile"
	"github.com/acoshift/hime"
	redisstore "github.com/acoshift/session/store/redis"
	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"

	"github.com/acoshift/himetic/app"
)

func main() {
	config := configfile.NewReader("config")

	db, err := sql.Open("postgres", config.String("db_url"))
	if err != nil {
		log.Fatalf("main: open database error; %v", err)
	}

	sessionStorage := redisstore.New(redisstore.Config{
		Prefix: config.String("session_prefix"),
		Pool: &redis.Pool{
			IdleTimeout: time.Minute,
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", config.String("session_host"))
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Now().Sub(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
		},
	})

	log.Println("main: start web server on :8080")
	err = hime.New().
		Handler(app.Handler(app.Config{
			SessionStorage: sessionStorage,
			SessionName:    config.String("session_name"),
			SessionSecret:  config.Bytes("session_secret"),
			DB:             db,
		})).
		GracefulShutdown().
		ListenAndServe(":8080")
	if err != nil {
		log.Fatalf("main: start web server error; %v", err)
	}
}
