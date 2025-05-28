package usecases

import (
	"context"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
	"lanchonete/usecases"
)

type confirmarPagamentoUseCase struct {
	pagamentoRepository repository.PagamentoRepository
}

func NewConfirmarPagamentoUseCase(pagamentoRepository repository.PagamentoRepository) usecases.ConfirmarPagamentoUseCase {
	return &confirmarPagamentoUseCase{
		pagamentoRepository: pagamentoRepository,
	}
}

func (uc *confirmarPagamentoUseCase) ConfirmarPagamento(c context.Context, pagamento *entities.Pagamento) error {
	return uc.pagamentoRepository.ConfirmarPagamento(c, pagamento)
}

