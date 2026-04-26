package controllers

import (
	"net/http"
	"webapp/src/cookies"
)

// FazerLogout remove os dados de autenticação salvos no browser do usuário
func FazerLogout(res http.ResponseWriter, req *http.Request) {
	cookies.Deletar(res)
	http.Redirect(res, req, "/login", 302)
}
