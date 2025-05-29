package repository

import (
	"context"
	"lanchonete/internal/domain/entities"
)

type AcompanhamentoRepository interface {
	CriarAcompanhamento(c context.Context) (int, error)
	AdicionarPedido(c context.Context, idAcompanhamento int, idPedido int) error
	AtualizarStatusPedido(c context.Context, idPedido int, novoStatus entities.StatusPedido) error
	BuscarAcompanhamento(c context.Context, idAcompanhamento int) (*entities.AcompanhamentoPedido, error)
	BuscarPedidos(ctx context.Context, idPedido int) ([]entities.Pedido, error)

}
