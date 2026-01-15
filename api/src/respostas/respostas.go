package respostas

import (
	"encoding/json"
	"log"
	"net/http"
)

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

// Erro retorna um erro em formato Json
func Erro(res http.ResponseWriter, statusCode int, erro error) {
	JSON(res, statusCode, struct {
		Erro string `json:"erro"`
	}{
		Erro: erro.Error(),
	})
}
