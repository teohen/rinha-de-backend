package pessoa

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/teohen/rinha-de-backend/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, pessoa domain.Pessoa) (error, uuid.UUID)
	Test(ctx context.Context)
}

type pessoaRepository struct {
	db *pgxpool.Pool
}

func NewPessoaRepository(db *pgxpool.Pool) Repository {
	return &pessoaRepository{
		db: db,
	}
}

func (p *pessoaRepository) Create(ctx context.Context, pessoa domain.Pessoa) (error, uuid.UUID) {
	insert := "INSERT INTO pessoas (apelido, uid, nome, nascimento, stack) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING returnin uid;"

	err := p.db.QueryRow(ctx, insert, pessoa.Apelido, pessoa.UUID, pessoa.Nome, pessoa.Nascimento, pessoa.Stack).Scan(&pessoa.UUID)

	if err != nil {
		return fmt.Errorf("create pessoa: %w", err), uuid.Nil
	}

	return nil, pessoa.UUID
}

func (p *pessoaRepository) Test(ctx context.Context) {
	var seg int

	err := p.db.QueryRow(ctx, "SELECT 1 + 2 AS result;").Scan(&seg)

	if err != nil {
		fmt.Println("ERRO NO SQL", err)
	}

	fmt.Println("SEG", seg)
}
