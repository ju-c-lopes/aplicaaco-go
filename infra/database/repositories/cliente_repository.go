package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
)

type clienteMysqlRepository struct {
	db *sql.DB
}

func NewClienteMysqlRepository(db *sql.DB) repository.ClienteRepository {
	return &clienteMysqlRepository{
		db: db,
	}
}

func (cr *clienteMysqlRepository) CriarCliente(c context.Context, cliente *entities.Cliente) error {
	query := "INSERT INTO Cliente (cpfCliente, nomeCliente, emailCliente) VALUES (?, ?, ?)"

	_, err := cr.db.ExecContext(c, query, cliente.CPF, cliente.Nome, cliente.Email)
	if err != nil {
		return fmt.Errorf("erro ao criar cliente: %w", err)
	}
	return nil
}

func (cr *clienteMysqlRepository) BuscarCliente(c context.Context, CPF string) (*entities.Cliente, error) {
	query := "SELECT cpfCliente, nomeCliente, emailCliente FROM Cliente WHERE cpfCliente = ?"

	var cliente entities.Cliente
	err := cr.db.QueryRowContext(c, query, CPF).Scan(&cliente.CPF, &cliente.Nome, &cliente.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cliente não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar cliente: %w", err)
	}
	return &cliente, nil
}
