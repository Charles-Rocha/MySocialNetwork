package controllers

import (
	"net/http"
	"webapp/src/utils"
)

// CarregarTelaDeLogin irá carregar a tela de login
func CarregarTelaDeLogin(res http.ResponseWriter, req *http.Request) {
	utils.ExecutarTemplate(res, "login.html", nil)
}

// CarregarPaginaDeCadastroDoUsuario irá carregar a página de cadastro do usuário
func CarregarPaginaDeCadastroDoUsuario(res http.ResponseWriter, req *http.Request) {
	utils.ExecutarTemplate(res, "cadastro.html", nil)
}
