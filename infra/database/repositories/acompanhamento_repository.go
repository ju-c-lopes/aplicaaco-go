package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
)

type acompanhamentoMySQLRepository struct {
	db *sql.DB
}

func NewAcompanhamentoMySQLRepository(db *sql.DB) repository.AcompanhamentoRepository {
	return &acompanhamentoMySQLRepository{db: db}
}

func (r *acompanhamentoMySQLRepository) CriarAcompanhamento(ctx context.Context) (int, error) {
	result, err := r.db.ExecContext(ctx, `INSERT INTO Acompanhamento (tempoEstimado) VALUES ('00:15:00')`)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return int(id), nil
}

func (r *acompanhamentoMySQLRepository) AdicionarPedido(ctx context.Context, idAcompanhamento int, idPedido int) error {
	var ordem int
	_ = r.db.QueryRowContext(ctx,
		`SELECT IFNULL(MAX(ordem), 0) + 1 FROM FilaPedidos WHERE idAcompanhamento = ?`, idAcompanhamento).
		Scan(&ordem)

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO FilaPedidos (idAcompanhamento, idPedido, ordem) VALUES (?, ?, ?)`,
		idAcompanhamento, idPedido, ordem)
	return err
}

func (r *acompanhamentoMySQLRepository) AtualizarStatusPedido(ctx context.Context, idPedido int, novoStatus entities.StatusPedido) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE Pedido SET status = ? WHERE idPedido = ?`, string(novoStatus), idPedido)
	if err != nil {
		return err
	}

	if novoStatus == entities.Finalizado {
		_, err := r.db.ExecContext(ctx, `DELETE FROM FilaPedidos WHERE idPedido = ?`, idPedido)
		return err
	}
	return nil
}

func (r *acompanhamentoMySQLRepository) BuscarFila(ctx context.Context, idAcompanhamento int) (*entities.AcompanhamentoPedido, error) {
	query := `
	SELECT p.idPedido, p.status, p.totalPedido
	FROM FilaPedidos f
	JOIN Pedido p ON f.idPedido = p.idPedido
	WHERE f.idAcompanhamento = ?
	ORDER BY f.ordem ASC`

	rows, err := r.db.QueryContext(ctx, query, idAcompanhamento)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar fila: %v", err)
	}
	defer rows.Close()

	var pedidos []entities.Pedido
	for rows.Next() {
		var p entities.Pedido
		err := rows.Scan(&p.ID, &p.Status, &p.Total)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear pedido: %v", err)
		}
		pedidos = append(pedidos, p)
	}

	var ultimaAtualizacao time.Time
	var tempoEstimado string
	err = r.db.QueryRowContext(ctx,
		`SELECT tempoEstimado, ultimaAtualizacao FROM Acompanhamento WHERE idAcompanhamento = ?`, idAcompanhamento).
		Scan(&tempoEstimado, &ultimaAtualizacao)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar metadados do acompanhamento: %v", err)
	}

	return &entities.AcompanhamentoPedido{
		ID:                idAcompanhamento,
		Pedidos:           pedidos,
		TempoEstimado:     tempoEstimado,
		UltimaAtualizacao: ultimaAtualizacao,
	}, nil
}

// BuscarAcompanhamento implements repository.AcompanhamentoRepository.
func (r *acompanhamentoMySQLRepository) BuscarAcompanhamento(ctx context.Context, idAcompanhamento int) (*entities.AcompanhamentoPedido, error) {
	var ultimaAtualizacao time.Time
	var tempoEstimado string
	err := r.db.QueryRowContext(ctx,
		`SELECT tempoEstimado, ultimaAtualizacao FROM Acompanhamento WHERE idAcompanhamento = ?`, idAcompanhamento).
		Scan(&tempoEstimado, &ultimaAtualizacao)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar acompanhamento: %v", err)
	}

	return &entities.AcompanhamentoPedido{
		ID:                idAcompanhamento,
		TempoEstimado:     tempoEstimado,
		UltimaAtualizacao: ultimaAtualizacao,
	}, nil
}

func (r *acompanhamentoMySQLRepository) BuscarPedidos(ctx context.Context, idPedido int) ([]entities.Pedido, error) {
	query := `
		SELECT 
			p.idPedido, 
			p.cliente, 
			c.nomeCliente, 
			p.totalPedido, 
			p.tempoEstimado, 
			p.status, 
			p.statusPagamento
		FROM FilaPedidos f
		JOIN Pedido p ON f.idPedido = p.idPedido
		JOIN Cliente c ON p.cliente = c.cpfCliente
		WHERE f.idAcompanhamento = ?
		ORDER BY f.ordem;
	`

	rows, err := r.db.QueryContext(ctx, query, idPedido)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pedidos da fila: %w", err)
	}
	defer rows.Close()

	var pedidos []entities.Pedido

	for rows.Next() {
		var p entities.Pedido
		var tempoEstimado string
		var clienteCPF string
		var clienteNome string

		if err := rows.Scan(
			&p.ID,
			&clienteCPF,
			&clienteNome,
			&p.Total,
			&tempoEstimado,
			&p.Status,
			&p.StatusPagamento,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear pedido: %w", err)
		}

		p.ClienteCPF = fmt.Sprintf("%s - %s", clienteCPF, clienteNome)
		p.Produtos = []entities.Produto{}

		// Buscar produtos
		prodQuery := `SELECT pr.idProduto, pr.nomeProduto, pr.descricaoProduto, pr.precoProduto, pr.personalizacaoProduto, pr.categoriaProduto FROM Pedido_Produto pp JOIN Produto pr ON pp.idProduto = pr.idProduto WHERE pp.idPedido = ?`
		prodRows, err := r.db.QueryContext(ctx, prodQuery, p.ID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar produtos do pedido: %w", err)
		}

		for prodRows.Next() {
			var prod entities.Produto
			var personalizacao sql.NullString

			if err := prodRows.Scan(
				&prod.ID,
				&prod.Nome,
				&prod.Descricao,
				&prod.Preco,
				&personalizacao,
				&prod.Categoria,
			); err != nil {
				prodRows.Close()
				return nil, fmt.Errorf("erro ao escanear produto: %w", err)
			}

			prod.Personalizacao = personalizacao
			p.Produtos = append(p.Produtos, prod)
		}
		prodRows.Close()

		pedidos = append(pedidos, p)
	}

	return pedidos, nil
}

