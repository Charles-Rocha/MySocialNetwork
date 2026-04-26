package modelos

import "time"

//Publicacao representa uma publicação feita pelo usuário
type Publicacao struct {
	ID          uint64    `json:"id,omitempty"`
	Titulo      string    `json:"titulo,omitempty"`
	Conteudo    string    `json:"conteudo,omitempty"`
	AutorUserId uint64    `json:"autoruserId,omitempty"`
	AutorNick   string    `json:"autorNick,omitempty"`
	Curtidas    uint64    `json:"curtidas"`
	CriadaEm    time.Time `json:"criadaEm,omitempty"`
}
