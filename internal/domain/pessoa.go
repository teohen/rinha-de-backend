package domain

type Pessoa struct {
	UUID       uuid.UUID `json:"id"`
	Apelido    string    `json:"apelido"`
	Nome       string    `json:"Nome"`
	Nascimento string    `json:"nascimento"`
	Stack      []string  `json:"stack"`
}
