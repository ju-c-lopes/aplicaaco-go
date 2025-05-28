package bootstrap

import (
	"context"
	"database/sql"
	"log"

	"lanchonete/infra/database"
	"lanchonete/internal/domain/repository"
)

type App struct {
	Env                    *Env
	DB                       *sql.DB
	AcompanhamentoRepository repository.AcompanhamentoRepository
	PedidoRepository         repository.PedidoRepository
	ProdutoRepository        repository.ProdutoRepository
	ClienteRepository        repository.ClienteRepository
	PagamentoRepository      repository.PagamentoRepository
}

func NewApp(ctx context.Context) (*App, error) {
	// Load environment variables
	env := NewEnv()

	db, err := database.NewMySQLConnection(
		env.DBUser,
		env.DBPass,
		env.DBHost,
		env.DBPort,
		env.DBName,
	)

	if err != nil {
		log.Fatalf("erro ao conectar ao MySQL: %v", err)
	}

	// Initialize repositories
	acompanhamentoRepo, pedidoRepo, produtoRepo, clienteRepo, pagamentoRepo := NewRepositories(db)

	return &App{
		Env:                    env,
		DB:                     db,
		AcompanhamentoRepository: acompanhamentoRepo,
		PedidoRepository:         pedidoRepo,
		ProdutoRepository:        produtoRepo,
		ClienteRepository:        clienteRepo,
		PagamentoRepository:      pagamentoRepo,
	}, nil
}