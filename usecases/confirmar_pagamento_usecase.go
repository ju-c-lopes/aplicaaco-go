package usecases

import (
	"context"
	"lanchonete/internal/domain/entities"
)

type ConfirmarPagamentoUseCase interface {
	ConfirmarPagamento(C context.Context, pagamento *entities.Pagamento) error
}