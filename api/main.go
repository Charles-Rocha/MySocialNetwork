package main

import (
	"api/src/config"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

// Função usada apenas para criar a chave secreta para poder salvar no arquivo .env
/* func init() {
	chave := make([]byte, 64)
	if _, erro := rand.Read(chave); erro != nil {
		log.Fatal(erro)
	}

	stringBase64 := base64.StdEncoding.EncodeToString(chave)
	fmt.Println(stringBase64)
} */

func main() {
	config.Carregar()
	//fmt.Println(config.StringConexaoBanco)
	//fmt.Println(config.Porta)

	fmt.Printf("Servidor escutando na porta %d\n", config.Porta)
	r := router.Gerar()

	//fmt.Println(config.SecretKey)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), r))
}
