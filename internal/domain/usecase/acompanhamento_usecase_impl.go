package usecase

import (
	"context"
	"fmt"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
)

type acompanhamentoUseCase struct {
	acompanhamentoRepo repository.AcompanhamentoRepository
}

func NewAcompanhamentoUseCase(repo repository.AcompanhamentoRepository) AcompanhamentoUseCase {
	return &acompanhamentoUseCase{acompanhamentoRepo: repo}
}

func (uc *acompanhamentoUseCase) CriarAcompanhamento(ctx context.Context) (int, error) {
	return uc.acompanhamentoRepo.CriarAcompanhamento(ctx)
}

func (uc *acompanhamentoUseCase) AdicionarPedido(ctx context.Context, idAcompanhamento int, idPedido int) error {
	return uc.acompanhamentoRepo.AdicionarPedido(ctx, idAcompanhamento, idPedido)
}

func (uc *acompanhamentoUseCase) BuscarAcompanhamento(ctx context.Context, idAcompanhamento int) (*entities.AcompanhamentoPedido, error) {
	return uc.acompanhamentoRepo.BuscarAcompanhamento(ctx, idAcompanhamento)
}

func (uc *acompanhamentoUseCase) AtualizarStatusPedido(ctx context.Context, idPedido int, novoStatus entities.StatusPedido) error {
	fmt.Printf("UseCase: atualizando pedido %d para status %s\n", idPedido, novoStatus)
	return uc.acompanhamentoRepo.AtualizarStatusPedido(ctx, idPedido, novoStatus)
}

func (uc *acompanhamentoUseCase) BuscarPedidos(ctx context.Context, idAcompanhamento int) ([]entities.Pedido, error) {
	return uc.acompanhamentoRepo.BuscarPedidos(ctx, idAcompanhamento)
}
