package usecases

import (
	"context"
	"errors"
	"lanchonete/internal/domain/entities"
	"testing"
)

// MockClienteRepository simulates the repository for Cliente
type MockClienteRepository struct {
	clientes []*entities.Cliente
}

func (m *MockClienteRepository) CriarCliente(ctx context.Context, cliente *entities.Cliente) error {
	for _, c := range m.clientes {
		if c.CPF == cliente.CPF {
			return errors.New("cliente já existe")
		}
	}
	m.clientes = append(m.clientes, cliente)
	return nil
}

func (m *MockClienteRepository) BuscarCliente(ctx context.Context, cpf string) (*entities.Cliente, error) {
	for _, c := range m.clientes {
		if c.CPF == cpf {
			return c, nil
		}
	}
	return &entities.Cliente{}, errors.New("cliente não encontrado")
}

// ClienteUseCaseImpl implements ClienteUseCase for testing
type ClienteUseCaseImpl struct {
	repo *MockClienteRepository
}

func (uc *ClienteUseCaseImpl) CriarCliente(c context.Context, cliente *entities.Cliente) error {
	return uc.repo.CriarCliente(c, cliente)
}

func (uc *ClienteUseCaseImpl) BuscarCliente(c context.Context, cpf string) (entities.Cliente, error) {
	cliente, err := uc.repo.BuscarCliente(c, cpf)
	if err != nil {
		return entities.Cliente{}, err
	}
	return *cliente, nil
}

func TestClienteUseCase_CriarEBuscarCliente(t *testing.T) {
	mockRepo := &MockClienteRepository{}
	useCase := &ClienteUseCaseImpl{repo: mockRepo}

	cliente, err := entities.ClienteNew("Ana", "ana@email.com", "12345678900")
	if err != nil {
		t.Fatalf("erro ao criar cliente: %v", err)
	}

	// Test creating the cliente
	err = useCase.CriarCliente(context.Background(), cliente)
	if err != nil {
		t.Fatalf("erro inesperado ao criar cliente: %v", err)
	}

	// Test fetching the cliente
	found, err := useCase.BuscarCliente(context.Background(), cliente.CPF)
	if err != nil {
		t.Fatalf("erro inesperado ao buscar cliente: %v", err)
	}
	if found.Nome != cliente.Nome || found.CPF != cliente.CPF || found.Email != cliente.Email {
		t.Errorf("atributos do cliente não conferem: obtido %+v, esperado %+v", found, cliente)
	}
}
