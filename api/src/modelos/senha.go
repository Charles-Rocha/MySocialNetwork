package modelos

//Senha representa o formato da requisição de alteração de senha
type Senha struct {
	Atual string `json:"atual"`
	Nova  string `json:"nova"`
}
