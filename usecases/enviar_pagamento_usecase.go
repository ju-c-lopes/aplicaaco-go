package usecases

import (
	"context"
	"lanchonete/internal/domain/entities"
)

type EnviarPagamentoUseCase interface {
	EnviarPagamento(c context.Context, pagamento *entities.Pagamento) error
}