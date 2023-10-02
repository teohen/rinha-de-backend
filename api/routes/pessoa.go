package routes

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/teohen/rinha-de-backend/api/handler"
	"github.com/teohen/rinha-de-backend/internal/pessoa"
)

type ServerAPI struct {
	Server *http.Server
}

func NewServer(conn *pgxpool.Pool) *ServerAPI {

	port := os.Getenv("HTTP_PORT")

	router := chi.NewRouter()

	repository := pessoa.NewPessoaRepository(conn)

	service := pessoa.NewService(repository)

	handler := handler.NewPessoaHandler(service)

	router.Post("/pessoas", handler.Create)
	router.Get("/pessoas/{id}", handler.Get)
	router.Get("/pessoas", handler.Search)
	router.Get("/contagem-pessoas", handler.Count)

	api := &ServerAPI{
		Server: &http.Server{
			Handler: router,
			Addr:    ":" + port,
		},
	}

	return api
}
