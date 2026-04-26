package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"webapp/src/config"
	"webapp/src/requisicoes"
	"webapp/src/respostas"

	"github.com/gorilla/mux"
)

// CriarPublicacao chama a API para cadastrar uma publicação no banco de dados
func CriarPublicacao(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("Criando publicação.")
	req.ParseForm()

	publicacao, erro := json.Marshal(map[string]string{
		"titulo":   req.FormValue("titulo"),
		"conteudo": req.FormValue("conteudo"),
	})

	if erro != nil {
		respostas.JSON(res, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	url := fmt.Sprintf("%s/publicacoes", config.ApiUrl)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(req, http.MethodPost, url, bytes.NewBuffer(publicacao))

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

// CurtirPublicacao chama a API para curtir uma publicação
func CurtirPublicacao(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	publicacaoId, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.JSON(res, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	url := fmt.Sprintf("%s/publicacoes/%d/curtir", config.ApiUrl, publicacaoId)
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

// DescurtirPublicacao chama a API para descurtir uma publicação
func DescurtirPublicacao(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	publicacaoId, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.JSON(res, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	url := fmt.Sprintf("%s/publicacoes/%d/descurtir", config.ApiUrl, publicacaoId)
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

// AtualizarPublicacao chama a API para atualizar uma publicação
func AtualizarPublicacao(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	publicacaoId, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.JSON(res, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	req.ParseForm()

	publicacao, erro := json.Marshal(map[string]string{
		"titulo":   req.FormValue("titulo"),
		"conteudo": req.FormValue("conteudo"),
	})

	if erro != nil {
		respostas.JSON(res, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	url := fmt.Sprintf("%s/publicacoes/%d", config.ApiUrl, publicacaoId)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(req, http.MethodPut, url, bytes.NewBuffer(publicacao))

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

// DeletarPublicacao chama a API para deletar uma publicação
func DeletarPublicacao(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	publicacaoId, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.JSON(res, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	url := fmt.Sprintf("%s/publicacoes/%d", config.ApiUrl, publicacaoId)
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
