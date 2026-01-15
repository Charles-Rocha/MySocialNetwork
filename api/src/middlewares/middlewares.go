package middlewares

import (
	"api/src/autenticacao"
	"api/src/respostas"
	"log"
	"net/http"
)

// Logger escreve informações da requisição no terminal
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("\n %s %s %s", req.Method, req.RequestURI, req.Host)
		next(res, req) //A função next executa a função que vem por parâmetro (next aqui no caso foi um nome qualquer que ele deu, poderia ser ProximaFuncao por exemplo)
	}
}

// Autenticar verifica se o usuário que está fazendo a requisição está autenticado
func Autenticar(next http.HandlerFunc) http.HandlerFunc { //http.HandlerFunc é igual a func (w http.ResponseWrite, r *http.Request)
	return func(res http.ResponseWriter, req *http.Request) {
		if erro := autenticacao.ValidarToken(req); erro != nil {
			respostas.Erro(res, http.StatusUnauthorized, erro)
			return
		}
		next(res, req) //A função next executa a função que vem por parâmetro
	}
}
