package pessoa

import "context"

type Service interface {
	Test(ctx context.Context)
}

type pessoaService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &pessoaService{
		repository: r,
	}
}

func (p *pessoaService) Test(ctx context.Context) {
	p.repository.Test(ctx)
}
