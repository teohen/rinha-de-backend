package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (api *serverAPi) HandlerPostPessoa(w http.ResponseWriter, r *http.Request) {

	fmt.Println("1")
	type param struct {
		Apelido    string   `json:"apelido"`
		Nome       string   `json:"nome"`
		Nascimento string   `json:"nascimento"`
		Stack      []string `json:"stack"`
	}
	fmt.Println("2")

	decoder := json.NewDecoder(r.Body)

	params := param{}
	err := decoder.Decode(&params)

	if err != nil {
		w.WriteHeader(400)
	}
	fmt.Println("3")

	pessoa := Pessoa{
		Apelido:    params.Apelido,
		Nome:       params.Nome,
		Nascimento: params.Nascimento,
		Stack:      params.Stack,
	}

	err, id := api.repo.create(context.Background(), pessoa)
	fmt.Println("4")

	if err != nil {
		fmt.Println("creating new pessoa: %w", err)
		respondWithError(w, 400, "Bad Request")
	}
	fmt.Println("5")

	respondWithJSON(w, 200, id)
}

func (api *serverAPi) test(w http.ResponseWriter, r *http.Request) {

	// TODO: corrigir o problema da conexão com o banco
	var teteo int
	api.repo.dbPool.QueryRow(context.Background(), "SELECT t1 FROM (VALUES ('teteo')) t1 (c2)").Scan(&teteo)

	fmt.Println("sfkdsjfskfj", teteo)
	w.WriteHeader(303)
}
