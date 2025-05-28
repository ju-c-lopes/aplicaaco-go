package presenters

import (
	"fmt"
	"lanchonete/internal/domain/entities"
	"time"
)

// AcompanhamentoDTO representa os dados de um acompanhamento para apresentação
type AcompanhamentoDTO struct {
	ID            string      `json:"id"`
	Pedidos       []PedidoDTO `json:"pedidos"`
	TempoEstimado int         `json:"tempoEstimado"` // in minutes
}

// NewAcompanhamentoDTO cria um novo DTO a partir de uma entidade AcompanhamentoPedido
func NewAcompanhamentoDTO(a *entities.AcompanhamentoPedido) *AcompanhamentoDTO {
	pedidos := make([]PedidoDTO, 0)
	for _, p := range a.Pedidos {
		pedidos = append(pedidos, *NewPedidoDTO(&p))
	}
	return &AcompanhamentoDTO{
		ID:            fmt.Sprintf("%d", a.ID),
		Pedidos:       pedidos,
		TempoEstimado: int(time.Until(a.TempoEstimado).Minutes()),
	}
}
