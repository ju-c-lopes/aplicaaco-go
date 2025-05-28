package usecases

import (
	"context"
	"errors"
	"lanchonete/internal/domain/entities"
	"testing"
	"time"
)

// MockPedidoRepository implements repository.PedidoRepository for testing
type MockPedidoRepository struct {
	Pedidos []*entities.Pedido
}

func (m *MockPedidoRepository) CriarPedido(ctx context.Context, pedido *entities.Pedido) error {
	// Simulate duplicate check
	for _, p := range m.Pedidos {
		if p.ID == pedido.ID {
			return errors.New("pedido já existe")
		}
	}
	m.Pedidos = append(m.Pedidos, pedido)
	return nil
}

func (m *MockPedidoRepository) BuscarPedido(ctx context.Context, id int) (*entities.Pedido, error) {
	for _, p := range m.Pedidos {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("pedido não encontrado")
}

func (m *MockPedidoRepository) AtualizarStatusPedido(ctx context.Context, pedidoID int, status string, ultimaAtualizacao time.Time) error {
	for _, p := range m.Pedidos {
		if p.ID == pedidoID {
			p.Status = entities.StatusPedido(status)
			p.UltimaAtualizacao = ultimaAtualizacao
			return nil
		}
	}
	return errors.New("pedido não encontrado")
}

func (m *MockPedidoRepository) ListarTodosOsPedidos(ctx context.Context) ([]*entities.Pedido, error) {
	return m.Pedidos, nil
}

func TestPedidoUseCase_Run_MultiplePedidos(t *testing.T) {
	mockRepo := &MockPedidoRepository{}
	useCase := NewPedidoIncluirUseCase(mockRepo)

	// Produtos base (same as in produto test)
	produtos := []entities.Produto{
		{Nome: "Hamburguer", Categoria: entities.Lanche, Descricao: "Hamburguer artesanal", Preco: 25.0},
		{Nome: "Batata Frita", Categoria: entities.Acompanhamento, Descricao: "Batata frita crocante", Preco: 10.0},
		{Nome: "Refrigerante", Categoria: entities.Bebida, Descricao: "Coca-Cola lata", Preco: 7.5},
	}

	pedidos := []struct {
		Cliente        entities.Cliente
		Produtos       []entities.Produto
	}{
		{
			Cliente:        entities.Cliente{Nome: "João", CPF: "11111111111"},
			Produtos:       []entities.Produto{produtos[0], produtos[1]},
		},
		{
			Cliente:        entities.Cliente{Nome: "Maria", CPF: "22222222222"},
			Produtos:       []entities.Produto{produtos[0], produtos[2]},
		},
		{
			Cliente:        entities.Cliente{Nome: "Pedro", CPF: "33333333333"},
			Produtos:       []entities.Produto{produtos[0], produtos[1], produtos[2]},
		},
	}

	for _, p := range pedidos {
		id, err := useCase.Run(context.Background(), p.Cliente.CPF, p.Produtos)
		if err != nil {
			t.Fatalf("unexpected error for pedido %+v: %v", p, err)
		}
		if id == nil {
			t.Fatalf("expected pedido to be created for %+v", p)
		}
	}

	if len(mockRepo.Pedidos) != 3 {
		t.Errorf("expected 3 pedidos in repository, got %d", len(mockRepo.Pedidos))
	}

	// Optionally, check attributes of each created pedido
	for i, pedido := range mockRepo.Pedidos {
		expected := pedidos[i]
		if pedido.ClienteCPF != expected.Cliente.CPF {
			t.Errorf("pedido cliente mismatch: got %+v, want %+v", pedido.ClienteCPF, expected.Cliente.CPF)
		}
		if len(pedido.Produtos) != len(expected.Produtos) {
			t.Errorf("pedido produtos count mismatch: got %d, want %d", len(pedido.Produtos), len(expected.Produtos))
		}
	}
}
