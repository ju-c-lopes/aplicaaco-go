package usecases

import (
	"context"
	"fmt"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
	"lanchonete/internal/domain/usecase"
)

type acompanhamentoUseCase struct {
	acompanhamentoRepo repository.AcompanhamentoRepository
}

func NewAcompanhamentoUseCase(acompanhamentoRepo repository.AcompanhamentoRepository) usecase.AcompanhamentoUseCase {
	return &acompanhamentoUseCase{
		acompanhamentoRepo: acompanhamentoRepo,
	}
}

func (uc *acompanhamentoUseCase) CriarAcompanhamento(c context.Context) (int, error) {
	return uc.acompanhamentoRepo.CriarAcompanhamento(c)
}


func (uc *acompanhamentoUseCase) AdicionarPedido(c context.Context, acompanhamentoID int, pedidoID int) error {
	return uc.acompanhamentoRepo.AdicionarPedido(c, acompanhamentoID, pedidoID)
}

func (uc *acompanhamentoUseCase) AtualizarStatusPedido(c context.Context, identificacao int, novoStatus entities.StatusPedido) error {
	fmt.Printf("UseCase: Atualizando pedido %d para status %s\n", identificacao, novoStatus)

	return uc.acompanhamentoRepo.AtualizarStatusPedido(c, identificacao, novoStatus)
}

// Implementa o método BuscarAcompanhamento para satisfazer a interface usecase.AcompanhamentoUseCase
func (uc *acompanhamentoUseCase) BuscarAcompanhamento(c context.Context, acompanhamentoID int) (*entities.AcompanhamentoPedido, error) {
	return uc.acompanhamentoRepo.BuscarAcompanhamento(c, acompanhamentoID)
}
