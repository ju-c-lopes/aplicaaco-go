package usecases

import (
	"context"
	"lanchonete/internal/domain/entities"
)

type AcompanhamentoCriarUseCase interface {
	CriarAcompanhamento(c context.Context, acompanhamento *entities.AcompanhamentoPedido) error
}
