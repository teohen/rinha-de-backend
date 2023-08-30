package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/teohen/rinha-de-backend/internal/domain"
	"github.com/teohen/rinha-de-backend/internal/pessoa"
)

type PessoaRequest struct {
	Apelido    string   `json:"apelido"`
	Nome       string   `json:"nome"`
	Nascimento string   `json:"nascimento"`
	Stack      []string `json:"stack"`
}

type Handler interface {
	Test(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
	Count(w http.ResponseWriter, r *http.Request)
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

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	newPessoa := domain.Pessoa{
		Apelido:    pessoaParams.Apelido,
		Nome:       pessoaParams.Nome,
		Nascimento: pessoaParams.Nascimento,
		Stack:      pessoaParams.Stack,
	}

	valid := validate(newPessoa)

	if valid == false {
		respondWithError(w, http.StatusUnprocessableEntity, "unprocessable entity")
		return
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

	w.Header().Add("Location", fmt.Sprintf("/pessoas/%s", pessoaUid))
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

	if pessoa.Nome == "" {
		respondWithError(w, 404, "Not Found")
		return
	}

	respondWithJSON(w, 200, pessoa)
	return
}

func (phandler *pessoaHandler) Search(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("t")

	if term == "" {
		respondWithError(w, 400, "t query string is required")
		return
	}

	err, pessoaList := phandler.service.Search(context.Background(), term)

	if pessoaList == nil {
		pessoaList = []domain.Pessoa{}
	}

	if err != nil {
		fmt.Println("Error:", err)
		respondWithError(w, 500, "internal server error")
		return
	}

	respondWithJSON(w, 200, pessoaList)
	return
}

func (phandler *pessoaHandler) Count(w http.ResponseWriter, r *http.Request) {

	err, count := phandler.service.Count(context.Background())
	if err != nil {
		fmt.Println("Error no count", err)
		respondWithError(w, 500, "internal server error")
		return
	}

	respondWithJSON(w, 200, count)
	return
}

func (phandler *pessoaHandler) Test(w http.ResponseWriter, r *http.Request) {

	phandler.service.Test(context.Background())

	w.WriteHeader(303)
}

func validate(pessoa domain.Pessoa) bool {
	validations := []bool{
		pessoa.Nome != "",
		len([]rune(pessoa.Nome)) <= 100,
		pessoa.Apelido != "",
		len([]rune(pessoa.Apelido)) <= 32,
		pessoa.Nascimento != ""}

	valid := checkValidations(validations)

	if valid == false {
		return false
	}

	_, err := time.Parse("2006-01-02", pessoa.Nascimento)

	if err != nil {
		return false
	}

	valid = validateStack(pessoa.Stack)

	if valid == false {
		return false
	}
	return true
}

func checkValidations(validations []bool) bool {
	valid := true
	for _, validation := range validations {
		if validation == false {
			valid = false
		}
	}
	return valid
}

func validateStack(stack []string) bool {
	valid := true
	if len(stack) > 0 {
		for _, item := range stack {
			validations := []bool{item != "", len(item) <= 32}
			if checkValidations(validations) == false {
				valid = false
				break
			}

		}
	}
	return valid
}
