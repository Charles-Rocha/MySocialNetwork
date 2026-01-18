package modelos

//DadosAutenticacao contém o token e o Id do usuário autenticado
type DadosAutenticacao struct {
	ID    string `json:"id"`
	Token string `jason:"token"`
}
