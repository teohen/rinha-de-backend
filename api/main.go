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

	var psqlconn string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s", os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	poolConfig, err := pgxpool.ParseConfig(psqlconn)

	if err != nil {
		log.Fatal("Couldn parse config", err)
	}

	poolConfig.MaxConns = 100
	poolConfig.MinConns = 10

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)

	if err != nil {
		log.Fatal("Couldnt create conn pool", err)
	}

	defer db.Close()

	server := routes.NewServer(db)

	log.Printf("Server running")

	err = server.Server.ListenAndServe()

	if err != nil {
		log.Fatal("Couldnt create server: ", err)
	}
}
