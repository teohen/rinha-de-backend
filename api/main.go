package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/rueidis"
	"github.com/teohen/rinha-de-backend/api/routes"
)

func main() {

	godotenv.Load()

	var psqlconn string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	poolConfig, err := pgxpool.ParseConfig(psqlconn)

	if err != nil {
		log.Fatal("Couldn parse config", err)
	}

	poolConfig.MaxConns = 150
	poolConfig.MinConns = 100

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)

	if err != nil {
		log.Fatal("Couldnt create conn pool", err)
	}

	redisClient, err := rueidis.NewClient(
		rueidis.ClientOption{InitAddress: []string{os.Getenv("REDIS_HOST") + ":6379"}},
	)

	if err != nil {
		log.Fatal("Couldnt connect with redis", err)
	}

	defer db.Close()

	defer redisClient.Close()

	server := routes.NewServer(db, redisClient)

	log.Printf("Server running")

	err = server.Server.ListenAndServe()

	if err != nil {
		log.Fatal("Couldnt create server: ", err)
	}
}
