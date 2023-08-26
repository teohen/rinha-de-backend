package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pessoa struct {
	UUID       uuid.UUID `json:"id"`
	Apelido    string    `json:"apelido"`
	Nome       string    `json:"Nome"`
	Nascimento string    `json:"nascimento"`
	Stack      []string  `json:"stack"`
}

type PessoaDB interface {
	create(ctx context.Context, pessoa Pessoa) (error, uuid.UUID)
	get(ctx context.Context, id uuid.UUID) (error, Pessoa)
	find(ctx context.Context, query string) ([]Pessoa, error)
}

type PessoaModel struct {
	store PessoaDB
}

type pessoaDBImpl struct {
	dbPool *pgxpool.Pool
}

func (p *pessoaDBImpl) create(ctx context.Context, pessoa Pessoa) (error, uuid.UUID) {
	insert := "INSERT INTO pessoas (apelido, uid, nome, nascimento, stack) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING returnin uid;"

	err := p.dbPool.QueryRow(ctx, insert, pessoa.Apelido, pessoa.UUID, pessoa.Nome, pessoa.Nascimento, pessoa.Stack).Scan(&pessoa.UUID)

	if err != nil {
		return fmt.Errorf("create pessoa: %w", err), uuid.Nil
	}

	return nil, pessoa.UUID
}

func (p *pessoaDBImpl) get(ctx context.Context, id uuid.UUID) (error, Pessoa) {
	var pessoa Pessoa
	get := "SELECT apelido, uid, nome, nascimento, stack FROM pessoas WHERE uid = $1"

	err := p.dbPool.QueryRow(ctx, get, id).Scan(&pessoa.Apelido, &pessoa.UUID, &pessoa.Nome, &pessoa.Nascimento, &pessoa.Stack)

	if err != nil {
		return fmt.Errorf("get pessoa: %w", err), pessoa
	}

	return nil, pessoa
}

func (p *pessoaDBImpl) find(ctx context.Context, query string) (error, []Pessoa) {
	var pessoas []Pessoa

	find := "SELECT apelido, uid, nome, nascimento, stack FROM pessoas ilike $1 limit 50"

	rows, err := p.dbPool.Query(ctx, find, "%"+query+"%")

	if err != nil {
		return fmt.Errorf("find pessoas: %w", err), pessoas
	}

	defer rows.Close()

	for rows.Next() {
		var pessoa Pessoa
		err := rows.Scan(&pessoa.Apelido, &pessoa.UUID, &pessoa.Nome, &pessoa.Nascimento, &pessoa.Stack)

		if err != nil {
			return fmt.Errorf("Scanning pessoa: %w", err), pessoas
		}

		pessoas = append(pessoas, pessoa)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("for loop over pessoas: %w", err), pessoas
	}

	return nil, pessoas
}
