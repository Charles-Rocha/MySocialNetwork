package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	ApiUrl   = ""   //ApiUrl representa a URL para comunicação com a API
	Porta    = 0    //Porta onde a aplicação web está rodando
	HashKey  []byte //HashKey é utilizada para autenticar o cookie
	BlockKey []byte //BlockKey é utilizada para criptografar os dados do cookie
)

// Carregar inicializa as variáveis de ambiente
func Carregar() {
	var erro error

	//Load simplesmente vai ler o arquivo .env (No Delphi a gente usa arquivo .ini (só para nível de comparação))
	if erro = godotenv.Load(); erro != nil {
		log.Fatal(erro)
	}

	Porta, erro = strconv.Atoi(os.Getenv("APP_PORT"))
	if erro != nil {
		log.Fatal(erro)
	}

	ApiUrl = os.Getenv("API_URL")
	HashKey = []byte(os.Getenv("HASH_KEY"))
	BlockKey = []byte(os.Getenv("BLOCK_KEY"))
}
