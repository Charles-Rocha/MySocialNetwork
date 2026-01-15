package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"fmt"
)

// usuarios representa um repositório de usuários
type usuarios struct {
	db *sql.DB
}

// NovoRepositorioDeUsuarios cria um repositório de usuários
func NovoRepositorioDeUsuarios(db *sql.DB) *usuarios {
	return &usuarios{db}
}

// Criar insere um usuário no banco de dados - Método
func (repositorio usuarios) Criar(usuario modelos.Usuario) (int, error) {
	statement, erro := repositorio.db.Prepare(
		"INSERT INTO usuarios1 (nome, nick, email, senha) VALUES(?, ?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)
	if erro != nil {
		return 0, erro
	}

	ultimoIdInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}
	return int(ultimoIdInserido), nil
}

// Listar traz todos os usuários que atendam a um filtro com nome ou nick
func (repositorio usuarios) Listar(nomeOuNick string) ([]modelos.Usuario, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick) //%nomeOuNick%

	linhas, erro := repositorio.db.Query(
		"SELECT id, nome, nick, email, criadoem from usuarios1 "+
			"WHERE nome LIKE ? or nick LIKE ?", nomeOuNick, nomeOuNick)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []modelos.Usuario
	for linhas.Next() {
		var usuario modelos.Usuario
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, usuario)
	}
	return usuarios, nil
}

// ListarUsuarioPorId traz o usuário do banco de dados da tabela usuarios através do seu Id
func (repositorio usuarios) ListarUsuarioPorId(usuarioId int64) (modelos.Usuario, error) {
	linhas, erro := repositorio.db.Query(
		"SELECT id, nome, nick, email, criadoem from usuarios1 WHERE Id = ?", usuarioId,
	)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linhas.Close()

	var usuario modelos.Usuario
	if linhas.Next() {
		if erro = linhas.Scan(
			&usuarioId,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return modelos.Usuario{}, erro
		}
	}
	return usuario, nil
}

// AtualizarUsuario altera as informações de um usuário do banco de dados da tabela usuarios através do seu Id
func (repositorio usuarios) AtualizarUsuario(usuarioId uint64, usuario modelos.Usuario) error {
	statement, erro := repositorio.db.Prepare(
		"UPDATE usuarios1 SET nome=?, nick=?, email=? WHERE Id=?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuarioId); erro != nil {
		return erro
	}

	return nil
}

// DeletarUsuario deletar um usuário do banco de dados da tabela usuarios através do seu Id
func (repositorio usuarios) DeletarUsuario(usuarioId uint64) error {
	statement, erro := repositorio.db.Prepare(
		"DELETE FROM usuarios1 WHERE Id=?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(usuarioId); erro != nil {
		return erro
	}

	return nil
}

// ListarUsuarioPorEmail traz o usuário do banco de dados da tabela usuarios através do seu Email
func (repositorio usuarios) ListarUsuarioPorEmail(usuarioEmail string) (modelos.Usuario, error) {
	linha, erro := repositorio.db.Query("SELECT id, senha FROM usuarios1 WHERE email = ?", usuarioEmail)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linha.Close()

	var usuario modelos.Usuario
	if linha.Next() {
		if erro = linha.Scan(&usuario.ID, &usuario.Senha); erro != nil {
			return modelos.Usuario{}, erro
		}
	}

	return usuario, nil
}

// Seguir permite que um usuário siga outro
func (repositorio usuarios) Seguir(usuarioId, seguidorId uint64) error {
	statement, erro := repositorio.db.Prepare(
		//IGNORE simplesmente ignora a inserção de uma chave duplicada. Se eu já tiver por exemplo que o id 1 já está seguindo o id 2 e as mesmas são chaves compostas
		//na tabela seguidores1, então não posso permitir que seja inserido novamente. Com o IGNORE ele vai ignorar a inserção duplicada
		"INSERT IGNORE INTO seguidores1(usuario_id, seguidor_id) VALUES(?, ?)",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(usuarioId, seguidorId); erro != nil {
		return erro
	}

	return nil
}

// PararDeSeguir permite que um usuário pare de seguir outro
func (repositorio usuarios) PararDeSeguir(usuarioId, seguidorId uint64) error {
	statement, erro := repositorio.db.Prepare(
		"DELETE FROM seguidores1 where usuario_id = ? AND seguidor_id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(usuarioId, seguidorId); erro != nil {
		return erro
	}

	return nil
}

// ListarSeguidores lista todos os seguidores de um usuario
func (repositorio usuarios) ListarSeguidores(usuarioId uint64) ([]modelos.Usuario, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT u.id, u.nome, u.nick, u.email, u.criadoem FROM usuarios1 u
		INNER JOIN seguidores1 s on s.seguidor_id = u.id
		WHERE s.usuario_id = ?`, usuarioId)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []modelos.Usuario
	for linhas.Next() {
		var usuario modelos.Usuario

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

// ListarSeguindo lista todos os usuário que um determinado usuário está seguindo
func (repositorio usuarios) ListarSeguindo(usuarioId uint64) ([]modelos.Usuario, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT u.id, u.nome, u.nick, u.email, u.criadoem FROM usuarios1 u
		INNER JOIN seguidores1 s on s.usuario_id = u.id
		WHERE s.seguidor_id = ?`, usuarioId)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []modelos.Usuario
	for linhas.Next() {
		var usuario modelos.Usuario

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

// ListarSenha lista a senha de um usuario através de seu Id
func (repositorio usuarios) ListarSenha(usuarioId uint64) (string, error) {
	linha, erro := repositorio.db.Query("SELECT senha FROM usuarios1 WHERE id = ?", usuarioId)
	if erro != nil {
		return "", erro
	}
	defer linha.Close()

	var usuario modelos.Usuario

	if linha.Next() {
		if erro = linha.Scan(&usuario.Senha); erro != nil {
			return "", erro
		}
	}

	return usuario.Senha, nil
}

// AtualizarSenha altera a senha de um usuário
func (repositorio usuarios) AtualizarSenha(usuarioId uint64, senha string) error {
	statement, erro := repositorio.db.Prepare("UPDATE usuarios1 set senha = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(senha, usuarioId); erro != nil {
		return erro
	}

	return nil
}
