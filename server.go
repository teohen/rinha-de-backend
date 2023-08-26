package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type serverAPi struct {
	server *http.Server
	repo   *pessoaDBImpl
}

func NewServer(conn *pgxpool.Pool) *serverAPi {

	port := os.Getenv("HTTP_PORT")

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	pessoaDB := pessoaDBImpl{
		dbPool: conn,
	}

	api := &serverAPi{
		server: &http.Server{
			Handler: router,
			Addr:    ":" + port,
		},
		repo: &pessoaDB,
	}

	router.Post("/pessoas", api.HandlerPostPessoa)
	router.Post("/test", api.test)

	return api

}
