package controllers

import (
	"net/http"
	"webapp/src/utils"
)

// CarregarTelaDeLogin irá carregar a tela de login
func CarregarTelaDeLogin(res http.ResponseWriter, req *http.Request) {
	utils.ExecutarTemplate(res, "login.html", nil)
}

// CarregarPaginaDeCadastroDeUsuario irá carregar a página de cadastro do usuário
func CarregarPaginaDeCadastroDeUsuario(res http.ResponseWriter, req *http.Request) {
	utils.ExecutarTemplate(res, "cadastro.html", nil)
}
