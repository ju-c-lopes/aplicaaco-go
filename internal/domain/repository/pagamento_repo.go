package repository

import ("context"
	"lanchonete/internal/domain/entities"
)

type PagamentoRepository interface {
	EnviarPagamento(c context.Context, pagamento *entities.Pagamento) error
	ConfirmarPagamento(c context.Context, pagamento *entities.Pagamento) error
}
