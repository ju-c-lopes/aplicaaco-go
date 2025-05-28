package entities

import "time"

type AcompanhamentoPedido struct {
	ID                int
	Pedidos           []Pedido
	TempoEstimado     time.Time
	UltimaAtualizacao time.Time
}
