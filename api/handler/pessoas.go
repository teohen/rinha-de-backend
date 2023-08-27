package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/teohen/rinha-de-backend/internal/pessoa"
)

type Handler interface {
	Test(w http.ResponseWriter, r *http.Request)
}

type pessoaHandler struct {
	service pessoa.Service
}

func NewPessoaHandler(pessoaService pessoa.Service) Handler {
	return &pessoaHandler{
		service: pessoaService,
	}
}

func (phandler *pessoaHandler) Test(w http.ResponseWriter, r *http.Request) {

	phandler.service.Test(context.Background())

	fmt.Println("sfkdsjfskfj")
	w.WriteHeader(303)
}
