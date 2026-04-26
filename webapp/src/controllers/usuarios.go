package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"webapp/src/config"
	"webapp/src/cookies"
	"webapp/src/requisicoes"
	"webapp/src/respostas"

	"github.com/gorilla/mux"
)

// CriarUsuario chama a API para cadastrar um usuário no banco de dados
func CriarUsuario(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	//fmt.Println(nome)

	usuario, erro := json.Marshal(map[string]string{
		"nome":  req.FormValue("nome"),
		"email": req.FormValue("email"),
		"nick":  req.FormValue("nick"),
		"senha": req.FormValue("senha"),
	})

	if erro != nil {
		respostas.JSON(res, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	//fmt.Println(bytes.NewBuffer(usuario))

	url := fmt.Sprintf("%s/usuarios", config.ApiUrl)
	response, erro := http.Post(url, "application/json", bytes.NewBuffer(usuario))
	if erro != nil {
		respostas.JSON(res, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(res, response)
		return
	}

	respostas.JSON(res, response.StatusCode, nil)

	//fmt.Println(response.Body)
}

// PararDeSeguirUsuario chama a api para parar de seguir um usuário
func PararDeSeguirUsuario(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	usuarioId, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.JSON(res, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	url := fmt.Sprintf("%s/usuarios/%d/parar-de-seguir", config.ApiUrl, usuarioId)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(req, http.MethodPost, url, nil)
	if erro != nil {
		respostas.JSON(res, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(res, response)
		return
	}

	respostas.JSON(res, response.StatusCode, nil)
}

// SeguirUsuario chama a api para seguir um usuário
func SeguirUsuario(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	usuarioId, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.JSON(res, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	url := fmt.Sprintf("%s/usuarios/%d/seguir", config.ApiUrl, usuarioId)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(req, http.MethodPost, url, nil)
	if erro != nil {
		respostas.JSON(res, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(res, response)
		return
	}

	respostas.JSON(res, response.StatusCode, nil)
}

// EditarUsuario chama a api para editar o usuário
func EditarUsuario(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	usuario, erro := json.Marshal(map[string]string{
		"nome":  req.FormValue("nome"),
		"nick":  req.FormValue("nick"),
		"email": req.FormValue("email"),
	})

	if erro != nil {
		respostas.JSON(res, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	cookie, _ := cookies.Ler(req)
	usuarioId, _ := strconv.ParseUint(cookie["id"], 10, 64)

	url := fmt.Sprintf("%s/usuarios/%d", config.ApiUrl, usuarioId)

	response, erro := requisicoes.FazerRequisicaoComAutenticacao(req, http.MethodPut, url, bytes.NewBuffer(usuario))
	if erro != nil {
		respostas.JSON(res, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(res, response)
		return
	}

	respostas.JSON(res, response.StatusCode, nil)
}

// AtualizarSenha chama a api para atualizar a senha do usuário
func AtualizarSenha(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	senhas, erro := json.Marshal(map[string]string{
		"atual": req.FormValue("atual"),
		"nova":  req.FormValue("nova"),
	})

	if erro != nil {
		respostas.JSON(res, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	cookie, _ := cookies.Ler(req)
	usuarioId, _ := strconv.ParseUint(cookie["id"], 10, 64)

	url := fmt.Sprintf("%s/usuarios/%d/atualizar-senha", config.ApiUrl, usuarioId)

	response, erro := requisicoes.FazerRequisicaoComAutenticacao(req, http.MethodPost, url, bytes.NewBuffer(senhas))
	if erro != nil {
		respostas.JSON(res, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(res, response)
		return
	}

	respostas.JSON(res, response.StatusCode, nil)
}

// DeletarUsuario chama a api para deletar o usuário
func DeletarUsuario(res http.ResponseWriter, req *http.Request) {
	cookie, _ := cookies.Ler(req)
	usuarioId, _ := strconv.ParseUint(cookie["id"], 10, 64)

	url := fmt.Sprintf("%s/usuarios/%d", config.ApiUrl, usuarioId)

	response, erro := requisicoes.FazerRequisicaoComAutenticacao(req, http.MethodDelete, url, nil)
	if erro != nil {
		respostas.JSON(res, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(res, response)
		return
	}

	respostas.JSON(res, response.StatusCode, nil)
}
