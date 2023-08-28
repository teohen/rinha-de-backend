package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/teohen/rinha-de-backend/internal/domain"
	"github.com/teohen/rinha-de-backend/internal/pessoa"
)

type PessoaRequest struct {
	Apelido    string   `json:"apelido" validate:"required, max=32"`
	Nome       string   `json:"nome" validate:"required, max=100"`
	Nascimento string   `json:"nascimento" validate:"required, datetime=2020-01-01"`
	Stack      []string `json:"stack" validate:"dive,max=32"`
}

type PessoaResponse struct {
	ID         string   `json:"id"`
	Apelido    string   `json:"apelido" validate:"required,max=32"`
	Nome       string   `json:"nome" validate:"required,max=100"`
	Nascimento string   `json:"nascimento" validate:"required,datetime=2006-01-02"`
	Stack      []string `json:"stack" validate:"dive,max=32"`
}

type Handler interface {
	Test(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
}

type pessoaHandler struct {
	service pessoa.Service
}

func NewPessoaHandler(pessoaService pessoa.Service) Handler {
	return &pessoaHandler{
		service: pessoaService,
	}
}

func (phandler *pessoaHandler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	pessoaParams := PessoaRequest{}

	err := decoder.Decode(&pessoaParams)

	newPessoa := domain.Pessoa{
		Apelido:    pessoaParams.Apelido,
		Nome:       pessoaParams.Nome,
		Nascimento: pessoaParams.Nascimento,
		Stack:      pessoaParams.Stack,
	}

	newPessoa.UUID = uuid.New()

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json: %s", err))
		return
	}

	err, pessoaUid := phandler.service.Create(context.Background(), newPessoa)

	if err != nil {
		fmt.Println("Error:", err)
		respondWithError(w, 500, "internal server error")
	}

	respondWithJSON(w, 201, pessoaUid)
}

func (phandler *pessoaHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid, err := uuid.Parse(id)

	if err != nil {
		respondWithError(w, 400, "id not formatted")
		return
	}

	err, pessoa := phandler.service.Get(context.Background(), uid)

	if err != nil {
		fmt.Println("Error:", err)
		respondWithError(w, 500, "internal server error")
		return
	}

	respondWithJSON(w, 200, pessoa)
	return
}

func (phandler *pessoaHandler) Search(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("t")

	err, pessoaList := phandler.service.Search(context.Background(), term)

	if err != nil {
		fmt.Println("Error:", err)
		respondWithError(w, 500, "internal server error")
		return
	}

	respondWithJSON(w, 200, pessoaList)
	return
}

func (phandler *pessoaHandler) Test(w http.ResponseWriter, r *http.Request) {

	phandler.service.Test(context.Background())

	fmt.Println("sfkdsjfskfj")
	w.WriteHeader(303)
}
