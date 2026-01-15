package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// CriarUsuario inseri um usuário no banco de dados na tabela usuarios
func CriarUsuario(res http.ResponseWriter, req *http.Request) {
	corpoRequest, erro := ioutil.ReadAll(req.Body)
	if erro != nil {
		respostas.Erro(res, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	if erro := usuario.Preparar("cadastro"); erro != nil {
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
	usuario.ID, erro = repositorio.Criar(usuario)
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(res, http.StatusCreated, usuario)
}

// ListarUsuarios lista todos os usuários salvos no banco de dados da tabela usuarios
func ListarUsuarios(res http.ResponseWriter, req *http.Request) {
	nomeOuNick := strings.ToLower(req.URL.Query().Get("usuario"))

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.Listar(nomeOuNick)
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(res, http.StatusOK, usuarios)
}

// ListarUsuario lista um usuário específico por Id salvo no banco de dados da tabela usuarios
func ListarUsuario(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)

	usuarioID, erro := strconv.ParseInt(parametros["usuarioId"], 10, 64)
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

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario, erro := repositorio.ListarUsuarioPorId(usuarioID)
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(res, http.StatusOK, usuario)
}

// AtualizarUsuario atualiza os dados de um usuário específico salvo no banco de dados da tabela usuarios
func AtualizarUsuario(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsarioId(req)
	if erro != nil {
		respostas.Erro(res, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Erro(res, http.StatusForbidden, errors.New("não é possível atualizar um usuário que não seja o seu"))
		return
	}

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

	if erro = usuario.Preparar("edicao"); erro != nil {
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
	erro = repositorio.AtualizarUsuario(usuarioID, usuario)
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(res, http.StatusNoContent, nil)
}

// DeletarUsuario remove um usuário específico do banco de dados da tabela usuarios
func DeletarUsuario(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsarioId(req)
	if erro != nil {
		respostas.Erro(res, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Erro(res, http.StatusForbidden, errors.New("não é possível excluir um usuário que não seja o seu"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	erro = repositorio.DeletarUsuario(usuarioID)
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(res, http.StatusNoContent, nil)
}

// SeguirUsuario permite que um usuário siga outro
func SeguirUsuario(res http.ResponseWriter, req *http.Request) {
	//seguidorId é o Id do usuário que está vindo na requisição, pela rota /usuarios/{usuarioId}/seguir
	seguidorId, erro := autenticacao.ExtrairUsarioId(req)
	if erro != nil {
		respostas.Erro(res, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(req)
	usuarioId, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	if seguidorId == usuarioId {
		respostas.Erro(res, http.StatusBadRequest, errors.New("não é possível seguir você mesmo"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Seguir(usuarioId, seguidorId); erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(res, http.StatusNoContent, nil)
}

// PararDeSeguirUsuario permite que um usuário para de deguir outro
func PararDeSeguirUsuario(res http.ResponseWriter, req *http.Request) {
	seguidorId, erro := autenticacao.ExtrairUsarioId(req)
	if erro != nil {
		respostas.Erro(res, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(req)
	usuarioId, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	if seguidorId == usuarioId {
		respostas.Erro(res, http.StatusForbidden, errors.New("não é possível parar de seguir a você mesmo"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.PararDeSeguir(usuarioId, seguidorId); erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(res, http.StatusNoContent, erro)
}

// ListarSeguidores lista todos os seguidores de um usuário
func ListarSeguidores(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	usuarioId, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
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

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	seguidores, erro := repositorio.ListarSeguidores(usuarioId)
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(res, http.StatusOK, seguidores)
}

// ListarSeguindo lista todos os usuários que um determinado usuário está seguindo
func ListarSeguindo(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	usuarioId, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
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

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.ListarSeguindo(usuarioId)
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(res, http.StatusOK, usuarios)
}

// AtualizarSenha permite atualizar a senha de um usuário
func AtualizarSenha(res http.ResponseWriter, req *http.Request) {
	usuarioIdNoToken, erro := autenticacao.ExtrairUsarioId(req)
	if erro != nil {
		respostas.Erro(res, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(req)
	usuarioId, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	if usuarioIdNoToken != usuarioId {
		respostas.Erro(res, http.StatusForbidden, errors.New("não é possível atualizar a senha de um usuário que não seja o seu"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(req.Body)
	var senha modelos.Senha
	if erro = json.Unmarshal(corpoRequisicao, &senha); erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	senhaSalvaNoBanco, erro := repositorio.ListarSenha(usuarioId)
	if erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	if erro = seguranca.VerificarSenha(senhaSalvaNoBanco, senha.Atual); erro != nil {
		respostas.Erro(res, http.StatusUnauthorized, errors.New("a senha atual não condiz com a senha salva no banco"))
		return
	}

	senhaComHash, erro := seguranca.Hash(senha.Nova)
	if erro != nil {
		respostas.Erro(res, http.StatusBadRequest, erro)
		return
	}

	if erro = repositorio.AtualizarSenha(usuarioId, string(senhaComHash)); erro != nil {
		respostas.Erro(res, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(res, http.StatusNoContent, nil)
}
