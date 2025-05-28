package bootstrap

import (
	"database/sql"
	"lanchonete/infra/database/repositories"
	"lanchonete/internal/domain/repository"
)

func NewRepositories(db *sql.DB) (
	acomp repository.AcompanhamentoRepository,
	pedido repository.PedidoRepository,
	produto repository.ProdutoRepository,
	cliente repository.ClienteRepository,
	pagamento repository.PagamentoRepository,
) {
	acomp = repositories.NewAcompanhamentoMySQLRepository(db)
	pedido = repositories.NewPedidoMysqlRepository(db)
	produto = repositories.NewProdutoMysqlRepository(db)
	cliente = repositories.NewClienteMysqlRepository(db)
	pagamento = repositories.NewPagamentoMysqlRepository(db)

	return
}
