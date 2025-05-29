package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
)

type pedidoMysqlRepository struct {
	db *sql.DB
}

func NewPedidoMysqlRepository(db *sql.DB) repository.PedidoRepository {
	return &pedidoMysqlRepository{db: db}
}

func (pr *pedidoMysqlRepository) CriarPedido(c context.Context, pedido *entities.Pedido) error {
	tx, err := pr.db.BeginTx(c, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	query := `INSERT INTO Pedido (cliente, totalPedido, tempoEstimado, status, statusPagamento) VALUES (?, ?, ?, ?, ?)`
	res, err := tx.ExecContext(c, query,
		pedido.ClienteCPF,
		pedido.Total,
		"00:15:00",
		pedido.Status,
		pedido.StatusPagamento,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao inserir pedido: %w", err)
	}

	pedidoID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao obter ID do pedido: %w", err)
	}
	pedido.ID = int(pedidoID)

	// Inserir produtos relacionados
	prodQuery := `INSERT INTO Pedido_Produto (idPedido, idProduto, quantidade) VALUES (?, ?, ?)`
	for _, prod := range pedido.Produtos {
		_, err := tx.ExecContext(c, prodQuery, pedidoID, prod.ID, 1)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("erro ao inserir produto no pedido: %w", err)
		}
	}

	return tx.Commit()
}

func (pr *pedidoMysqlRepository) BuscarPedido(c context.Context, identificacao int) (*entities.Pedido, error) {
	query := `SELECT idPedido, cliente, totalPedido, tempoEstimado, status, statusPagamento FROM Pedido WHERE idPedido = ?`

	var pedido entities.Pedido
	var clienteCPF string
	var tempoEstimado string

	fmt.Println("TimeStamp: ", pedido.TimeStamp, "ID: ", identificacao)
	err := pr.db.QueryRowContext(c, query, identificacao).Scan(
		&pedido.ID,
		&clienteCPF,
		&pedido.Total,
		&tempoEstimado,
		&pedido.Status,
		&pedido.StatusPagamento,
	)
	fmt.Println("Repository pedido: ", pedido.TimeStamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pedido não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar pedido: %w", err)
	}

	pedido.Produtos = []entities.Produto{}
	pedido.ClienteCPF = clienteCPF

	// Buscar produtos
	prodQuery := `SELECT idProduto FROM Pedido_Produto WHERE idPedido = ?`
	rows, err := pr.db.QueryContext(c, prodQuery, identificacao)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos do pedido: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p entities.Produto
		if err := rows.Scan(&p.Nome); err != nil {
			return nil, fmt.Errorf("erro ao escanear produto: %w", err)
		}
		pedido.Produtos = append(pedido.Produtos, p)
	}

	return &pedido, nil
}

func (pr *pedidoMysqlRepository) AtualizarStatusPedido(c context.Context, identificacao int, status string, ultimaAtualizacao time.Time) error {
	query := `UPDATE Pedido SET status = ?, ultimaAtualizacao = ? WHERE idPedido = ?`
	_, err := pr.db.ExecContext(c, query, status, ultimaAtualizacao, identificacao)
	if err != nil {
		return fmt.Errorf("erro ao atualizar status do pedido: %w", err)
	}
	return nil
}

func (pr *pedidoMysqlRepository) ListarTodosOsPedidos(c context.Context) ([]*entities.Pedido, error) {
	query := `SELECT idPedido, cliente, totalPedido, tempoEstimado, status, statusPagamento FROM Pedido`

	rows, err := pr.db.QueryContext(c, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pedidos: %w", err)
	}
	defer rows.Close()

	var pedidos []*entities.Pedido
	for rows.Next() {
		var p entities.Pedido
		var clienteCPF string
		var tempoEstimadoStr string

		if err := rows.Scan(
			&p.ID,
			&clienteCPF,
			&p.Total,
			&tempoEstimadoStr,
			&p.Status,
			&p.StatusPagamento,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear pedido: %w", err)
		}

		// pega a data atual (ano, mês e dia)
		now := time.Now()
		dataHoje := now.Format("2006-01-02")  // ex: "2025-05-28"

		// concatena uma data fixa com o horário recebido
		datetimeStr := dataHoje + " " + tempoEstimadoStr  // data fixa arbitrária

		// parse para time.Time completo
		_, err := time.Parse("2006-01-02 15:04:05", datetimeStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter tempoEstimado: %w", err)
		}
					
		p.TimeStamp = "00:15:00"  // Definindo um valor fixo para o TimeStamp
		p.ClienteCPF = clienteCPF
		p.Produtos = []entities.Produto{}

		// Buscar ids dos produtos do pedido
		prodIDsQuery := `SELECT idProduto FROM Pedido_Produto WHERE idPedido = ?`
		prodIDRows, err := pr.db.QueryContext(c, prodIDsQuery, p.ID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar ids dos produtos: %w", err)
		}

		var produtos []entities.Produto

		for prodIDRows.Next() {
			var idProduto int
			if err := prodIDRows.Scan(&idProduto); err != nil {
				prodIDRows.Close()
				return nil, fmt.Errorf("erro ao escanear idProduto: %w", err)
			}

			// Busca nomeProduto para esse idProduto
			nomeQuery := `SELECT nomeProduto FROM Produto WHERE idProduto = ?`
			var nomeProduto string
			err = pr.db.QueryRowContext(c, nomeQuery, idProduto).Scan(&nomeProduto)
			if err != nil {
				prodIDRows.Close()
				return nil, fmt.Errorf("erro ao buscar nomeProduto para id %d: %w", idProduto, err)
			}

			p := entities.Produto{
				ID: idProduto,
				Nome:   nomeProduto,
			}

			produtos = append(produtos, p)
		}
		prodIDRows.Close()

		if err := prodIDRows.Err(); err != nil {
			return nil, fmt.Errorf("erro na iteração dos ids dos produtos: %w", err)
		}

		// Atribui a lista de produtos ao pedido
		p.Produtos = produtos

		pedidos = append(pedidos, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro na iteração dos pedidos: %w", err)
	}

	return pedidos, nil
}
