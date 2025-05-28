package usecases

import (
	"context"
	"errors"
	"lanchonete/internal/domain/entities"
	"testing"
)

// MockProdutoRepository implements repository.ProdutoRepository for testing
type MockProdutoRepository struct {
	Produtos []*entities.Produto
}

func (m *MockProdutoRepository) AdicionarProduto(ctx context.Context, produto *entities.Produto) error {
	// Simulate duplicate check
	for _, p := range m.Produtos {
		if p.Nome == produto.Nome {
			return errors.New("produto já existe")
		}
	}
	m.Produtos = append(m.Produtos, produto)
	return nil
}

func (m *MockProdutoRepository) BuscarProdutoPorId(ctx context.Context, id int) (*entities.Produto, error) {
	for _, p := range m.Produtos {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("produto não encontrado")
}

func (m *MockProdutoRepository) ProdutoBuscarPorId(ctx context.Context, id int) (*entities.Produto, error) {
	for _, p := range m.Produtos {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("produto não encontrado")
}

func (m *MockProdutoRepository) EditarProduto(ctx context.Context, produto *entities.Produto) error {
	for i, p := range m.Produtos {
		if p.ID == produto.ID {
			m.Produtos[i] = produto
			return nil
		}
	}
	return errors.New("produto não encontrado")
}

func (m *MockProdutoRepository) RemoverProduto(ctx context.Context, id int) error {
	for i, p := range m.Produtos {
		if p.ID == id {
			m.Produtos = append(m.Produtos[:i], m.Produtos[i+1:]...)
			return nil
		}
	}
	return errors.New("produto não encontrado")
}

func (m *MockProdutoRepository) ListarTodosOsProdutos(ctx context.Context) ([]*entities.Produto, error) {
	return m.Produtos, nil
}

func (m *MockProdutoRepository) ListarPorCategoria(ctx context.Context, categoria string) ([]*entities.Produto, error) {
	var result []*entities.Produto
	for _, p := range m.Produtos {
		if string(p.Categoria) == categoria {
			result = append(result, p)
		}
	}
	return result, nil
}

func TestProdutoUseCase_Run_MultipleProducts(t *testing.T) {
	mockRepo := &MockProdutoRepository{}
	useCase := NewProdutoIncluirUseCase(mockRepo)

	produtos := []struct {
		Identificacao int
		Nome          string
		Categoria     string
		Descricao     string
		Preco         float32
	}{
		{1, "Hamburguer", "Lanche", "Hamburguer artesanal", 25.0},
		{2, "Batata Frita", "Acompanhamento", "Batata frita crocante", 10.0},
		{3, "Refrigerante", "Bebida", "Coca-Cola lata", 7.5},
	}

	for _, p := range produtos {
		produto, err := useCase.Run(context.Background(), p.Nome, p.Categoria, p.Descricao, p.Preco)
		if err != nil {
			t.Fatalf("unexpected error for produto %s: %v", p.Nome, err)
		}
		if produto == nil {
			t.Fatalf("expected produto to be created for %s", p.Nome)
		}
		if produto.Nome != p.Nome || string(produto.Categoria) != p.Categoria || produto.Descricao != p.Descricao || produto.Preco != p.Preco {
			t.Errorf("produto attributes mismatch: got %+v, want %+v", produto, p)
		}
	}

	if len(mockRepo.Produtos) != 3 {
		t.Errorf("expected 3 produtos in repository, got %d", len(mockRepo.Produtos))
	}
}
