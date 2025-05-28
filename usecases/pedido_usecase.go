package usecases

import (
	"context"
	"lanchonete/internal/domain/entities"
)

type PedidoUseCase interface {
	CriarPedido(c context.Context, pedido *entities.Pedido) error
	BuscarPedido(c context.Context, IdDoPedido int) (*entities.Pedido, error)
}
