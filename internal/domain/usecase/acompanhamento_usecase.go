package usecase

import (
	"context"
	"lanchonete/internal/domain/entities"
)

type AcompanhamentoUseCase interface {
	CriarAcompanhamento(ctx context.Context) (int, error)
	AdicionarPedido(ctx context.Context, idAcompanhamento int, idPedido int) error
	BuscarAcompanhamento(ctx context.Context, idAcompanhamento int) (*entities.AcompanhamentoPedido, error)
	AtualizarStatusPedido(ctx context.Context, idPedido int, novoStatus entities.StatusPedido) error
}
