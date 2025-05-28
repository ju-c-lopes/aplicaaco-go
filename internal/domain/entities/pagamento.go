package entities

import (
	"errors"
	"strings"
)

type Pagamento struct {
	IdPagamento  int
	IdPedido     int
	Valor        float64
	Status       string
	DataCriacao  string
}

func PagamentoNew(idPagamento, idPedido int, valor float64, status, dataCriacao string) (*Pagamento, error) {
	if idPagamento == 0 ||
		idPedido == 0 ||
		valor <= 0 ||
		strings.TrimSpace(status) == "" ||
		strings.TrimSpace(dataCriacao) == "" {
		return nil, errors.New("nenhum dos campos pode estar em branco")
	}

	return &Pagamento{
		IdPagamento: idPagamento,
		IdPedido:    idPedido,
		Valor:       valor,
		Status:      status,
		DataCriacao: dataCriacao,
	}, nil
}
