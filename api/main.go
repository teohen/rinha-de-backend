package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/teohen/rinha-de-backend/api/routes"
)

func main() {

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("Couldnt find DB_URL")
	}

	config, err := pgxpool.ParseConfig(dbURL)

	if err != nil {
		log.Fatal("Couldn parse config", err)
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		log.Fatal("Couldnt create conn pool", err)
	}

	defer dbPool.Close()

	server := NewServer(dbPool)

	log.Printf("Server running")

	err = server.server.ListenAndServe()

	if err != nil {
		log.Fatal("Couldnt create server: ", err)
	}

}
