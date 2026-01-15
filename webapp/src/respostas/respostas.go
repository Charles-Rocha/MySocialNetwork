package respostas

import (
	"encoding/json"
	"log"
	"net/http"
)

// Erro representa a reposta de erro da API (struct para mapear a mensagem de erro que está vindo na resposta da requisição para um Json)
type ErroAPI struct {
	Erro string `json:"erro"`
}

// JSON retorna uma resposta em Json para a requisição
func JSON(res http.ResponseWriter, statusCode int, dados interface{}) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	if dados != nil {
		if erro := json.NewEncoder(res).Encode(dados); erro != nil {
			log.Fatal(erro)
		}
	}
}

// TratarStatusCodeDeErro trata os StatusCode de erro das requisições (igual o maior que 400)
func TratarStatusCodeDeErro(res http.ResponseWriter, req *http.Response) {
	var erro ErroAPI
	json.NewDecoder(req.Body).Decode(&erro)
	JSON(res, req.StatusCode, erro)
}
