package usecases

import (
	"context"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
	"lanchonete/usecases"
)

type clienteUseCase struct {
	clienteRepository repository.ClienteRepository
}

func NewClienteUseCase(clienteRepository repository.ClienteRepository) usecases.ClienteUseCase {
	return &clienteUseCase{
		clienteRepository: clienteRepository,
	}
}

func (uc *clienteUseCase) CriarCliente(c context.Context, cliente *entities.Cliente) error {
	return uc.clienteRepository.CriarCliente(c, cliente)
}

func (uc *clienteUseCase) BuscarCliente(c context.Context, CPF string) (*entities.Cliente, error) {
	return uc.clienteRepository.BuscarCliente(c, CPF)
}