package autenticacao

import (
	"api/src/config"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// CriarToken retorna um token assinado com as permissões do usuário
func CriarToken(usuarioId uint64) (string, error) {
	permissoes := jwt.MapClaims{}
	permissoes["authorized"] = true
	permissoes["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissoes["usuarioId"] = usuarioId
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissoes) //Método de assinatura do token

	//É preciso assinar o token. Secret é achave que será utilizada para fazer a assinatura e garantir a autenticidade do token
	return token.SignedString([]byte(config.SecretKey))
}

// ValidarToken verifica se o token passado na requisição é válido
func ValidarToken(req *http.Request) error {
	tokenString := extrairToken(req)
	token, erro := jwt.Parse(tokenString, retornarChaveDeVerificacao)
	if erro != nil {
		return erro
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("token inválido")
}

// ExtrairUsarioId retorna o usuarioId que está salvo no token
func ExtrairUsarioId(req *http.Request) (uint64, error) {
	tokenString := extrairToken(req)
	token, erro := jwt.Parse(tokenString, retornarChaveDeVerificacao)
	if erro != nil {
		return 0, erro
	}

	if permissoes, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		usuarioID, erro := strconv.ParseUint(fmt.Sprintf("%.0f", permissoes["usuarioId"]), 10, 64)
		if erro != nil {
			return 0, erro
		}

		return usuarioID, nil
	}

	return 0, errors.New("token inválido")

}

func extrairToken(req *http.Request) string {
	token := req.Header.Get("Authorization")
	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return ""
}

func retornarChaveDeVerificacao(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("método de assinatura inesperado! %v", token.Header["alg"])
	}
	return config.SecretKey, nil
}
