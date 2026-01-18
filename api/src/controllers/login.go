package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Login é responsável por autenticar um usuário na API
func Login(res http.ResponseWriter, req *http.Request) {
	corpoRequisicao, erro := ioutil.ReadAll(req.Body)
	if erro != nil {
		respostas.Erro(res, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarioSalvoNoBanco, erro := repositorio.ListarUsuarioPorEmail(usuario.Email)
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	if erro = seguranca.VerificarSenha(usuarioSalvoNoBanco.Senha, usuario.Senha); erro != nil {
		respostas.Erro(res, http.StatusUnauthorized, erro)
		return
	}

	//res.Write([]byte("Autenticação feita com suceso!"))
	token, erro := autenticacao.CriarToken(usuarioSalvoNoBanco.ID)
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	usuarioId := strconv.FormatUint(usuarioSalvoNoBanco.ID, 10)

	respostas.JSON(res, http.StatusOK, modelos.DadosAutenticacao{ID: usuarioId, Token: token})

	fmt.Println(token)
	//res.Write([]byte(token)) //Resposta da requisição
}
