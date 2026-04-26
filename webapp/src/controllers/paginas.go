package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"webapp/src/config"
	"webapp/src/cookies"
	"webapp/src/modelos"
	"webapp/src/requisicoes"
	"webapp/src/respostas"
	"webapp/src/utils"

	"github.com/gorilla/mux"
)

// CarregarTelaDeLogin renderiza a tela de login
func CarregarTelaDeLogin(res http.ResponseWriter, req *http.Request) {
	/* cookie, _ := cookies.Ler(req)

	if cookie["token"] != "" {
		http.Redirect(res, req, "/home", 302)
		return
	} */
	utils.ExecutarTemplate(res, "login.html", nil)
}

// CarregarPaginaDeCadastroDeUsuario irá carregar a página de cadastro do usuário
func CarregarPaginaDeCadastroDeUsuario(res http.ResponseWriter, req *http.Request) {
	utils.ExecutarTemplate(res, "cadastro.html", nil)
}

// CarregarPaginaPrincipal irá carregar a página principal com as publicações
func CarregarPaginaPrincipal(res http.ResponseWriter, req *http.Request) {
	url := fmt.Sprintf("%s/publicacoes", config.ApiUrl)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(req, http.MethodGet, url, nil)

	if erro != nil {
		respostas.JSON(res, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(res, response)
		return
	}

	var publicacoes []modelos.Publicacao
	if erro = json.NewDecoder(response.Body).Decode(&publicacoes); erro != nil {
		respostas.JSON(res, http.StatusUnprocessableEntity, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	//fmt.Println(response.StatusCode, erro)

	// utils.ExecutarTemplate(res, "home.html", struct {
	// 	Publicacoes []modelos.Publicacao
	// 	OutroCampo  string
	// }{
	// 	Publicacoes: publicacoes,
	// 	OutroCampo:  "Valor Qualquer",
	// })

	cookie, _ := cookies.Ler(req)
	usuarioId, _ := strconv.ParseUint(cookie["id"], 10, 64)

	utils.ExecutarTemplate(res, "home.html", struct {
		Publicacoes []modelos.Publicacao
		UsuarioId   uint64
	}{
		Publicacoes: publicacoes,
		UsuarioId:   usuarioId,
	})
}

// CarregarPaginaDeAtualizacaoDePublicacao irá carregar a página de edição de publicação
func CarregarPaginaDeAtualizacaoDePublicacao(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	publicacaoId, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.JSON(res, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	url := fmt.Sprintf("%s/publicacoes/%d", config.ApiUrl, publicacaoId)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(req, http.MethodGet, url, nil)
	if erro != nil {
		respostas.JSON(res, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(res, response)
		return
	}

	var publicacao modelos.Publicacao
	if erro = json.NewDecoder(response.Body).Decode(&publicacao); erro != nil {
		respostas.JSON(res, http.StatusUnprocessableEntity, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	utils.ExecutarTemplate(res, "atualizar-publicacao.html", publicacao)
}

// CarregarPaginaDeUsuarios carrega a página com os usuários que atendam ao filtro passado
func CarregarPaginaDeUsuarios(res http.ResponseWriter, req *http.Request) {
	nomeOuNick := strings.ToLower(req.URL.Query().Get("usuario"))
	url := fmt.Sprintf("%s/usuarios?usuario=%s", config.ApiUrl, nomeOuNick)

	response, erro := requisicoes.FazerRequisicaoComAutenticacao(req, http.MethodGet, url, nil)
	if erro != nil {
		respostas.JSON(res, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(res, response)
		return
	}

	var usuarios []modelos.Usuario
	if erro := json.NewDecoder(response.Body).Decode(&usuarios); erro != nil {
		respostas.JSON(res, http.StatusUnprocessableEntity, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	utils.ExecutarTemplate(res, "usuarios.html", usuarios)
}

// CarregarPerfilDoUsuario irá carregar a página do perfil do usuário
func CarregarPerfilDoUsuario(res http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	usuarioId, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.JSON(res, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	usuario, erro := modelos.BuscarUsuarioCompleto(usuarioId, req)
	if erro != nil {
		respostas.JSON(res, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	cookie, _ := cookies.Ler(req)
	usuarioLogadoId, _ := strconv.ParseUint(cookie["id"], 10, 64)

	if usuarioId == usuarioLogadoId {
		http.Redirect(res, req, "/perfil", 302)
		return
	}

	utils.ExecutarTemplate(res, "usuario.html", struct {
		Usuario         modelos.Usuario
		UsuarioLogadoId uint64
	}{
		Usuario:         usuario,
		UsuarioLogadoId: usuarioLogadoId,
	})
}

// CarregarPerfilDoUsuarioLogado irá carregar a página do perfil do usuário logado
func CarregarPerfilDoUsuarioLogado(res http.ResponseWriter, req *http.Request) {
	cookie, _ := cookies.Ler(req)
	usuarioId, _ := strconv.ParseUint(cookie["id"], 10, 64)

	usuario, erro := modelos.BuscarUsuarioCompleto(usuarioId, req)
	if erro != nil {
		respostas.JSON(res, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	utils.ExecutarTemplate(res, "perfil.html", usuario)
}

// CarregarPaginaDeEdicaoDeUsuario irá carregar a página para edição dos dados do usuário
func CarregarPaginaDeEdicaoDeUsuario(res http.ResponseWriter, req *http.Request) {
	cookie, _ := cookies.Ler(req)
	usuarioId, _ := strconv.ParseUint(cookie["id"], 10, 64)

	canal := make(chan modelos.Usuario)
	go modelos.BuscarDadosDoUsuario(canal, usuarioId, req)
	usuario := <-canal

	if usuarioId == 0 {
		respostas.JSON(res, http.StatusInternalServerError, respostas.ErroAPI{Erro: "Erro ao buscar o usuário"})
		return
	}

	utils.ExecutarTemplate(res, "editar-usuario.html", usuario)
}

// CarregarPaginaDeAtualizacaoDeSenha irá carregar a página para atualização da senha do usuário
func CarregarPaginaDeAtualizacaoDeSenha(res http.ResponseWriter, req *http.Request) {
	utils.ExecutarTemplate(res, "atualizar-senha.html", nil)
}
