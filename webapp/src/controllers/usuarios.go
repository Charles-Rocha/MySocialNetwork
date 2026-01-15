package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"webapp/src/respostas"
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

	response, erro := http.Post("http://localhost:3300/usuarios", "application/json", bytes.NewBuffer(usuario))
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
