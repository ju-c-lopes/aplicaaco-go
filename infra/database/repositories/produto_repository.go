package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
)

type produtoMysqlRepository struct {
	database *sql.DB
}

func NewProdutoMysqlRepository(db *sql.DB) repository.ProdutoRepository {
	return &produtoMysqlRepository{
		database: db,
	}
}

func (pr *produtoMysqlRepository) AdicionarProduto(c context.Context, produto *entities.Produto) error {
	query := "INSERT INTO Produto (nomeProduto, descricaoProduto, precoProduto, personalizacaoProduto, categoriaProduto) VALUES (?, ?, ?, ?, ?)"
	_, err := pr.database.ExecContext(c, query, produto.Nome, produto.Descricao, produto.Preco, produto.Personalizacao, produto.Categoria)
	return err
}

func (pr *produtoMysqlRepository) BuscarProdutoPorId(c context.Context, id int) (*entities.Produto, error) {
	query := "SELECT idProduto, nomeProduto, descricaoProduto, precoProduto, personalizacaoProduto, categoriaProduto FROM Produto WHERE idProduto = ?"
	var produto entities.Produto
	var personalizacao sql.NullString
	err := pr.database.QueryRowContext(c, query, id).
		Scan(&produto.ID, &produto.Nome, &produto.Descricao, &produto.Preco, &personalizacao, &produto.Categoria)
	fmt.Println("Repository Buscando produto:", produto.Nome)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("produto não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar produto: %v", err)
	}
	if personalizacao.Valid {
		produto.Personalizacao = personalizacao
	} else {
		produto.Personalizacao = sql.NullString{Valid: false, String: ""}
	}
	fmt.Println("Repository Produto encontrado:", produto.Nome, produto.Descricao, produto.Preco, produto.Personalizacao, produto.Categoria)
	return &produto, nil
}

func (pr *produtoMysqlRepository) ListarTodosOsProdutos(c context.Context) ([]*entities.Produto, error) {
	query := "SELECT idProduto, nomeProduto, descricaoProduto, precoProduto, personalizacaoProduto, categoriaProduto FROM Produto"
	rows, err := pr.database.QueryContext(c, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos: %v", err)
	}
	defer rows.Close()

	var produtos []*entities.Produto
	var personalizacao sql.NullString
	for rows.Next() {
		var p entities.Produto
		var ident int
		if err := rows.Scan(&ident, &p.Nome, &p.Descricao, &p.Preco, &personalizacao, &p.Categoria); err != nil {
			return nil, fmt.Errorf("erro ao escanear produto: %v", err)
		}
		if personalizacao.Valid {
			p.Personalizacao = personalizacao
		} else {
			p.Personalizacao = sql.NullString{Valid: false}
		}

		produtos = append(produtos, &p)
	}
	return produtos, nil
}

func (pr *produtoMysqlRepository) EditarProduto(c context.Context, produto *entities.Produto) error {
	query := "UPDATE Produto SET nomeProduto = ?, descricaoProduto = ?, precoProduto = ?, personalizacaoProduto = ?, categoriaProduto = ? WHERE nomeProduto = ?"
	fmt.Println("Repository Atualizando produto:", produto.Nome, produto.Descricao, produto.Preco, produto.Personalizacao, produto.Categoria)
	result, err := pr.database.ExecContext(c, query, produto.Nome, produto.Descricao, produto.Preco, produto.Personalizacao, produto.Categoria, produto.Nome)
	if err != nil {
		return fmt.Errorf("erro ao atualizar produto: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar atualização: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("produto não encontrado")
	}

	return nil
}

func (pr *produtoMysqlRepository) RemoverProduto(c context.Context, id int) error {
	query := "DELETE FROM Produto WHERE idProduto = ?"
	result, err := pr.database.ExecContext(c, query, id)
	if err != nil {
		return fmt.Errorf("erro ao remover produto: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar remoção: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("produto não encontrado")
	}

	return nil
}

func (pr *produtoMysqlRepository) ListarPorCategoria(c context.Context, categoria string) ([]*entities.Produto, error) {
	query := "SELECT idProduto, nomeProduto, descricaoProduto, precoProduto, personalizacaoProduto, categoriaProduto FROM Produto WHERE categoriaProduto = ?"
	rows, err := pr.database.QueryContext(c, query, categoria)
	
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos por categoria: %v", err)
	}
	defer rows.Close()

	var produtos []*entities.Produto
	for rows.Next() {
		var p entities.Produto
		var ident int
		var personalizacao sql.NullString
		if err := rows.Scan(&ident, &p.Nome, &p.Descricao, &p.Preco, &personalizacao, &p.Categoria); err != nil {
			return nil, fmt.Errorf("erro ao escanear produto: %v", err)
		}
		if personalizacao.Valid {
			p.Personalizacao = personalizacao
		} else {
			p.Personalizacao = sql.NullString{Valid: false, String: ""}
		}
		produtos = append(produtos, &p)
	}
	return produtos, nil
}
