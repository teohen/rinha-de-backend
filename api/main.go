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

	var psqlconn string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	fmt.Println("SSS", psqlconn)
	poolConfig, err := pgxpool.ParseConfig(psqlconn)

	if err != nil {
		log.Fatal("Couldn parse config", err)
	}

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
