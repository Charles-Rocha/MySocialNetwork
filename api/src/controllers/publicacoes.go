package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CriarPublicacao adiciona uma nova publicação no banco de dados
func CriarPublicacao(res http.ResponseWriter, req *http.Request) {
	usuarioId, erro := autenticacao.ExtrairUsarioId(req)
	if erro != nil {
		respostas.Erro(res, http.StatusUnauthorized, erro)
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(req.Body)
	if erro != nil {
		respostas.Erro(res, http.StatusUnprocessableEntity, erro)
		return
	}

	var publicacao modelos.Publicacao
	if erro = json.Unmarshal(corpoRequisicao, &publicacao); erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	publicacao.AutorUserId = usuarioId

	if erro = publicacao.Preparar(); erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacao.ID, erro = repositorio.Criar(publicacao)
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(res, http.StatusCreated, publicacao)
}

// ListarPublicacoes traz as publicações que apareciam no feed do usuário
func ListarPublicacoes(res http.ResponseWriter, req *http.Request) {
	usuarioId, erro := autenticacao.ExtrairUsarioId(req)
	if erro != nil {
		respostas.Erro(res, http.StatusUnauthorized, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacoes, erro := repositorio.Listar(usuarioId)
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(res, http.StatusOK, publicacoes)
}

// ListarPublicacao traz uma única publicação
func ListarPublicacao(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	publicacaoId, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacao, erro := repositorio.ListarPorId(publicacaoId)
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(res, http.StatusOK, publicacao)
}

// AtualizarPublicacao altera os dados de uma publicação
func AtualizarPublicacao(res http.ResponseWriter, req *http.Request) {
	usuarioId, erro := autenticacao.ExtrairUsarioId(req)
	if erro != nil {
		respostas.Erro(res, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(req)
	publicacaoId, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacaoSalvaNoBanco, erro := repositorio.ListarPorId(publicacaoId)
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	if publicacaoSalvaNoBanco.AutorUserId != usuarioId {
		respostas.Erro(res, http.StatusForbidden, errors.New("não é possível atualizar uma publicação que não seja sua"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(req.Body)
	if erro != nil {
		respostas.Erro(res, http.StatusUnprocessableEntity, erro)
		return
	}

	var publicacao modelos.Publicacao
	if erro = json.Unmarshal(corpoRequisicao, &publicacao); erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	if erro = publicacao.Preparar(); erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	if erro = repositorio.Atualizar(publicacaoId, publicacao); erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(res, http.StatusNoContent, nil)

}

// DeletarPublicacao exclui os dados de uma publicação
func DeletarPublicacao(res http.ResponseWriter, req *http.Request) {
	usuarioId, erro := autenticacao.ExtrairUsarioId(req) //Lê o usuário ID que está no Token
	if erro != nil {
		respostas.Erro(res, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(req)                                                 //Pega os parâmetros da requisição
	publicacaoId, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64) //Pega o Id da publicação que está nos parâmetros da requisição
	if erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar() //Conecta com o Banco de Dados
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)          //Cria um novo repositório de publicação passando o banco de dados
	publicacaoSalvaNoBanco, erro := repositorio.ListarPorId(publicacaoId) //Pega uma publicação que está salva no banco passando o Id da publicação
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	//Verifico se o autor da publicação é o mesmo usuário que está fazendo a requisição, ou seja, se pertence a ele mesmo
	if publicacaoSalvaNoBanco.AutorUserId != usuarioId {
		respostas.Erro(res, http.StatusForbidden, errors.New("não é possível deletar uma publicação que não seja sua"))
		return
	}

	if erro = repositorio.Deletar(publicacaoId); erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(res, http.StatusNoContent, nil)
}

// ListarPublicacoesPorUsuario lista todas as publicações de um usuário especifício
func ListarPublicacoesPorUsuario(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)                                           //Pega os parâmetros da requisição
	usuarioId, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64) //Pega o Id da publicação que está nos parâmetros da requisição
	if erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar() //Conecta com o Banco de Dados
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db) //Cria um novo repositório de publicação passando o banco de dados
	publicacoes, erro := repositorio.ListarPorUsuario(usuarioId) //Pega uma publicação que está salva no banco passando o Id da publicação
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(res, http.StatusOK, publicacoes)
}

// CurtirPublicacao adiciona uma curtida na publicação
func CurtirPublicacao(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)                                                 //Pega os parâmetros da requisição
	publicacaoId, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64) //Pega o Id da publicação que está nos parâmetros da requisição
	if erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar() //Conecta com o Banco de Dados
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db) //Cria um novo repositório de publicação passando o banco de dados
	if erro = repositorio.Curtir(publicacaoId); erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(res, http.StatusNoContent, nil)
}

// DescurtirPublicacao remove uma curtida na publicação
func DescurtirPublicacao(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)                                                 //Pega os parâmetros da requisição
	publicacaoId, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64) //Pega o Id da publicação que está nos parâmetros da requisição
	if erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar() //Conecta com o Banco de Dados
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db) //Cria um novo repositório de publicação passando o banco de dados
	if erro = repositorio.Descurtir(publicacaoId); erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(res, http.StatusNoContent, nil)
}
