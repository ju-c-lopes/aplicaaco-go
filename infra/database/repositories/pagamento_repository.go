package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
)

type pagamentoMysqlRepository struct {
	db *sql.DB
}

func NewPagamentoMysqlRepository(db *sql.DB) repository.PagamentoRepository {
	return &pagamentoMysqlRepository{db: db}
}

func (pr *pagamentoMysqlRepository) EnviarPagamento(c context.Context, pagamento *entities.Pagamento) error {
	query := `INSERT INTO Pagamento (dataCriacao, Status, idPedido) VALUES (?, ?, ?)`

	dataCriacao := time.Now().Format("2006-01-02 15:04:05")

	res, err := pr.db.ExecContext(c, query,
		dataCriacao,
		pagamento.Status,
		pagamento.IdPedido,
	)
	if err != nil {
		return fmt.Errorf("erro ao enviar pagamento: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do pagamento: %w", err)
	}

	pagamento.IdPagamento = int(id)
	pagamento.DataCriacao = dataCriacao

	return nil
}

func (pr *pagamentoMysqlRepository) ConfirmarPagamento(c context.Context, pagamento *entities.Pagamento) error {
	query := `UPDATE Pagamento SET Status = ? WHERE idPagamento = ?`

	_, err := pr.db.ExecContext(c, query,
		pagamento.Status,
		pagamento.IdPagamento,
	)
	if err != nil {
		return fmt.Errorf("erro ao confirmar pagamento: %w", err)
	}

	return nil
}
