package pessoa

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/teohen/rinha-de-backend/internal/domain"
)

type Service interface {
	Create(ctx context.Context, pessoa domain.Pessoa) (error, uuid.UUID)
	Get(ctx context.Context, uid uuid.UUID) (error, domain.Pessoa)
	Search(ctx context.Context, term string) (error, []domain.Pessoa)
	Count(ctx context.Context) (error, int)
	GetByApelido(ctx context.Context, apelido string) (error, domain.Pessoa)
}

type pessoaService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &pessoaService{
		repository: r,
	}
}

func (p *pessoaService) Create(ctx context.Context, pessoa domain.Pessoa) (error, uuid.UUID) {
	err, pessoaRegistered := p.repository.GetByApelido(ctx, pessoa.Apelido)

	fmt.Println("pessoa registered", pessoaRegistered)
	if pessoaRegistered.UUID != uuid.Nil {

		return errors.New("pessoa already exists"), uuid.Nil
	}

	err, uidPessoa := p.repository.Create(ctx, pessoa)

	if err != nil {
		return err, uuid.Nil
	}

	return nil, uidPessoa
}

func (p *pessoaService) Get(ctx context.Context, uid uuid.UUID) (error, domain.Pessoa) {
	err, pessoa := p.repository.Get(ctx, uid)

	if err != nil {
		return err, domain.Pessoa{}
	}

	return nil, pessoa
}

func (p *pessoaService) GetByApelido(ctx context.Context, apelido string) (error, domain.Pessoa) {
	err, pessoa := p.repository.GetByApelido(ctx, apelido)

	if err != nil {
		return err, domain.Pessoa{}
	}

	return nil, pessoa
}

func (p *pessoaService) Search(ctx context.Context, term string) (error, []domain.Pessoa) {
	err, pessoaList := p.repository.Search(ctx, term)

	if err != nil {
		return err, nil
	}

	return nil, pessoaList

}

func (p *pessoaService) Count(ctx context.Context) (error, int) {
	err, count := p.repository.Count(ctx)

	if err != nil {
		return err, 0
	}

	return nil, count
}
