package middlewares

import (
	"log"
	"net/http"
	"webapp/src/cookies"
)

// Logger escreve informações da requisição no terminal
func Logger(proximaFuncao http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("\n %s %s %s", req.Method, req.RequestURI, req.Host)
		proximaFuncao(res, req)
	}
}

// Autenticar apenas verifica a existência de cookies, de qualquer forma, se os dados no cookie forem inválidos a api irá rejeitar
func Autenticar(proximaFuncao http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if _, erro := cookies.Ler(req); erro != nil {
			http.Redirect(res, req, "/login", 302)
			return
		}
		proximaFuncao(res, req)
	}
}
