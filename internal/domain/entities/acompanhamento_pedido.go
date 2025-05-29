package entities

import "time"

type AcompanhamentoPedido struct {
	ID                int         `json:"id,omitempty"`
	Pedidos           []Pedido    `json:"pedidos"`
	TempoEstimado     string      `json:"tempo_estimado"`
	UltimaAtualizacao time.Time   `json:"ultima_atualizacao"`
}
