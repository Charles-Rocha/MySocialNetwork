package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

// Publicacoes representa um repositório de publicações
type Publicacoes struct {
	db *sql.DB
}

// NovoRepositorioDePublicacoes cria um repositório de publicações
func NovoRepositorioDePublicacoes(db *sql.DB) *Publicacoes {
	return &Publicacoes{db}
}

// Criar insere uma publicação no banco de dados
func (repositorio Publicacoes) Criar(publicacao modelos.Publicacao) (uint64, error) {
	statement, erro := repositorio.db.Prepare(
		"INSERT INTO publicacoes1(titulo, conteudo, autor_user_id) VALUES(?, ?, ?) ",
	)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(publicacao.Titulo, publicacao.Conteudo, publicacao.AutorUserId)
	if erro != nil {
		return 0, erro
	}

	ultimoIdInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIdInserido), nil
}

// ListarPorId lista uma única publicação do banco de dados
func (repositorio Publicacoes) ListarPorId(publicacaoId uint64) (modelos.Publicacao, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT p.*, u.nick FROM publicacoes1 p
		INNER JOIN usuarios1 u on u.id = p.autor_user_id
		WHERE p.id = ?
	`, publicacaoId)

	if erro != nil {
		return modelos.Publicacao{}, erro
	}
	defer linhas.Close()

	var publicacao modelos.Publicacao
	if linhas.Next() {
		if erro = linhas.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorUserId,
			&publicacao.Curtidas,
			&publicacao.CriadaEm,
			&publicacao.AutorNick,
		); erro != nil {
			return modelos.Publicacao{}, erro
		}
	}

	return publicacao, nil
}

// Listar lista as publicações dos usuários seguidos e também do próprio usuário que fez a requisição
func (repositorio Publicacoes) Listar(usuarioId uint64) ([]modelos.Publicacao, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT DISTINCT p.*, u.nick FROM publicacoes1 p
		INNER JOIN usuarios1 u on u.id = p.autor_user_id
		INNER JOIN seguidores1 s on s.usuario_id = p.autor_user_id
		WHERE u.id = ? or s.seguidor_id = ?
		ORDER BY 1 DESC`, //ORDER BY 1 significa que ele vai ordenar pela primeira coluna dessa query, no caso aqui, o p.* que seria p.id
		usuarioId, usuarioId,
	)

	if erro != nil {
		return nil, erro
	}

	var publicacoes []modelos.Publicacao

	for linhas.Next() {
		var publicacao modelos.Publicacao

		if erro = linhas.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorUserId,
			&publicacao.Curtidas,
			&publicacao.CriadaEm,
			&publicacao.AutorNick,
		); erro != nil {
			return nil, erro
		}

		publicacoes = append(publicacoes, publicacao)
	}

	return publicacoes, nil
}

// Atualizar altera os dados de uma publicação no banco de dados
func (repositorio Publicacoes) Atualizar(publicacaoId uint64, publicacao modelos.Publicacao) error {
	statement, erro := repositorio.db.Prepare(`
		UPDATE publicacoes1 SET titulo = ?, conteudo = ? 
		WHERE id = ?
	`)

	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publicacao.Titulo, publicacao.Conteudo, publicacaoId); erro != nil {
		return erro
	}

	return nil
}

// Deletar exclui uma publicação no banco de dados
func (repositorio Publicacoes) Deletar(publicacaoId uint64) error {
	statement, erro := repositorio.db.Prepare(`
		DELETE FROM publicacoes1
		WHERE id = ?
	`)

	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publicacaoId); erro != nil {
		return erro
	}

	return nil
}

// ListarPorUsuario lista todas as publicações de um usuário específico
func (repositorio Publicacoes) ListarPorUsuario(usuarioId uint64) ([]modelos.Publicacao, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT p.*, u.nick FROM publicacoes1 p
		INNER JOIN usuarios1 u on u.id = p.autor_user_id		
		WHERE p.autor_user_id = ?`,
		usuarioId,
	)

	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var publicacoes []modelos.Publicacao

	for linhas.Next() {
		var publicacao modelos.Publicacao

		if erro = linhas.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorUserId,
			&publicacao.Curtidas,
			&publicacao.CriadaEm,
			&publicacao.AutorNick,
		); erro != nil {
			return nil, erro
		}

		publicacoes = append(publicacoes, publicacao)
	}

	return publicacoes, nil
}

// Curtir adiciona uma curtida na publicação
func (repositorio Publicacoes) Curtir(publicacaoId uint64) error {
	statement, erro := repositorio.db.Prepare(`
		UPDATE publicacoes1 set curtidas = curtidas + 1
		WHERE id = ?
	`)

	if erro != nil {
		return erro
	}

	defer statement.Close()

	if _, erro := statement.Exec(publicacaoId); erro != nil {
		return erro
	}

	return nil
}

// Descurtir remove uma curtida na publicação
func (repositorio Publicacoes) Descurtir(publicacaoId uint64) error {
	statement, erro := repositorio.db.Prepare(`
		UPDATE publicacoes1 set curtidas = 
		CASE 
			WHEN curtidas > 0 THEN curtidas -1
			ELSE 0
		END
		WHERE id = ?
	`)

	if erro != nil {
		return erro
	}

	defer statement.Close()

	if _, erro := statement.Exec(publicacaoId); erro != nil {
		return erro
	}

	return nil
}
