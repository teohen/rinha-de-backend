package pessoa

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/rueidis"
	"github.com/teohen/rinha-de-backend/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, pessoa domain.Pessoa) (error, uuid.UUID)
	Test(ctx context.Context)
	Get(ctx context.Context, id uuid.UUID) (error, domain.Pessoa)
	Search(ctx context.Context, term string) (error, []domain.Pessoa)
	Count(ctx context.Context) (error, int)
	SaveCache(ctx context.Context, pessoa domain.Pessoa) error
	GetCache(ctx context.Context, key string) (error, domain.Pessoa)
}

type pessoaRepository struct {
	db    *pgxpool.Pool
	cache rueidis.Client
}

func NewPessoaRepository(db *pgxpool.Pool, redis rueidis.Client) Repository {
	return &pessoaRepository{
		db:    db,
		cache: redis,
	}
}

func (p *pessoaRepository) Create(ctx context.Context, pessoa domain.Pessoa) (error, uuid.UUID) {

	stack := strings.Join(pessoa.Stack, ",")
	insert := "INSERT INTO pessoas (apelido, id, nome, nascimento, stack) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING returning id;"

	err := p.db.QueryRow(ctx, insert, pessoa.Apelido, pessoa.UUID, pessoa.Nome, pessoa.Nascimento, stack).Scan(&pessoa.UUID)

	if err != nil {
		return fmt.Errorf("create pessoa: %w", err), uuid.Nil
	}

	return nil, pessoa.UUID
}

func (p *pessoaRepository) Get(ctx context.Context, id uuid.UUID) (error, domain.Pessoa) {
	pessoa := domain.Pessoa{}

	var stack string
	get := "SELECT apelido, id, nome, nascimento, stack::varchar FROM pessoas WHERE id = $1"

	row := p.db.QueryRow(ctx, get, id)

	err := row.Scan(&pessoa.Apelido, &pessoa.UUID, &pessoa.Nome, &pessoa.Nascimento, &stack)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, domain.Pessoa{}
		}
		if err != nil {
			return fmt.Errorf("get pessoa: %w", err), pessoa
		}
	}

	pessoa.Stack = strings.Split(stack, ",")

	return nil, pessoa
}

func (p *pessoaRepository) Search(ctx context.Context, term string) (error, []domain.Pessoa) {
	query := `SELECT id, apelido, nome, nascimento, stack FROM pessoas p WHERE p.apelido ILIKE '%' || $1 || '%' LIMIT 50;`
	rows, err := p.db.Query(ctx, query, term)

	if err != nil {
		return err, nil
	}

	defer rows.Close()

	var pessoas []domain.Pessoa
	var stack string

	for rows.Next() {
		p := domain.Pessoa{}
		_ = rows.Scan(&p.UUID, &p.Apelido, &p.Nome, &p.Nascimento, &stack)
		p.Stack = strings.Split(stack, ",")
		pessoas = append(pessoas, p)
	}

	return nil, pessoas
}

func (p *pessoaRepository) Count(ctx context.Context) (error, int) {
	var count int

	err := p.db.QueryRow(ctx, "SELECT count(id) from pessoas").Scan(&count)

	if err != nil {
		fmt.Println("count", err)
		return err, 0
	}

	return nil, count
}

func (p *pessoaRepository) SaveCache(ctx context.Context, pessoa domain.Pessoa) error {
	pessoaMarshall, err := json.Marshal(pessoa)
	if err != nil {
		fmt.Println("Error marshalling pessoa struct", err)
		return err
	}

	pessoaString := string(pessoaMarshall)

	err = p.cache.Do(ctx, p.cache.B().Set().Key("pessoa:id:"+pessoa.UUID.String()).Value(pessoaString).Build()).Error()
	err = p.cache.Do(ctx, p.cache.B().Set().Key("pessoa:apelido:"+pessoa.Apelido).Value(pessoaString).Build()).Error()

	if err != nil {
		return err
	}

	return nil
}

func (p *pessoaRepository) GetCache(ctx context.Context, key string) (error, domain.Pessoa) {
	value, err := p.cache.DoCache(ctx, p.cache.B().Get().Key(key).Cache(), time.Minute).AsBytes()

	var pessoa domain.Pessoa

	err = json.Unmarshal(value, &pessoa)

	if err != nil {
		return err, domain.Pessoa{}
	}

	return nil, pessoa
}

func (p *pessoaRepository) Test(ctx context.Context) {
	var seg int

	err := p.db.QueryRow(ctx, "SELECT 1 + 2 AS result;").Scan(&seg)

	if err != nil {
		fmt.Println("ERRO NO SQL", err)
	}

	fmt.Println("SEG", seg)
}
