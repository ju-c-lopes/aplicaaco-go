package usecases

import (
	"context"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
	"lanchonete/usecases"
)

type enviarPagamentoUseCase struct {
	pagamentoRepository repository.PagamentoRepository
}

func NewEnviarPagamentoUseCase(pagamentoRepository repository.PagamentoRepository) usecases.EnviarPagamentoUseCase {
	return &enviarPagamentoUseCase{
		pagamentoRepository: pagamentoRepository,
	}
}

func (uc *enviarPagamentoUseCase) EnviarPagamento(c context.Context, pagamento *entities.Pagamento) error {
	return uc.pagamentoRepository.EnviarPagamento(c, pagamento)
}
