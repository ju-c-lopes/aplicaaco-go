package usecases

import (
	"context"
	"errors"
	"lanchonete/internal/domain/entities"
	"testing"
)

// MockPagamentoRepository simulates the repository for Pagamento
type MockPagamentoRepository struct {
	pagamentos []*entities.Pagamento
}

func (m *MockPagamentoRepository) EnviarPagamento(ctx context.Context, pagamento *entities.Pagamento) error {
	for _, p := range m.pagamentos {
		if p.IdPagamento == pagamento.IdPagamento {
			return errors.New("pagamento já existe")
		}
	}
	m.pagamentos = append(m.pagamentos, pagamento)
	return nil
}

func (m *MockPagamentoRepository) ConfirmarPagamento(ctx context.Context, pagamento *entities.Pagamento) error {
	for _, p := range m.pagamentos {
		if p.IdPagamento == pagamento.IdPagamento {
			p.Status = pagamento.Status
			return nil
		}
	}
	return errors.New("pagamento não encontrado")
}

// EnviarPagamentoUseCaseImpl implements EnviarPagamentoUseCase for testing
type EnviarPagamentoUseCaseImpl struct {
	repo *MockPagamentoRepository
}

func (uc *EnviarPagamentoUseCaseImpl) EnviarPagamento(ctx context.Context, pagamento *entities.Pagamento) error {
	return uc.repo.EnviarPagamento(ctx, pagamento)
}

// ConfirmarPagamentoUseCaseImpl implements ConfirmarPagamentoUseCase for testing
type ConfirmarPagamentoUseCaseImpl struct {
	repo *MockPagamentoRepository
}

func (uc *ConfirmarPagamentoUseCaseImpl) ConfirmarPagamento(ctx context.Context, pagamento *entities.Pagamento) error {
	return uc.repo.ConfirmarPagamento(ctx, pagamento)
}

func TestPagamentoUseCase_EnviarEConfirmarPagamento(t *testing.T) {
	mockRepo := &MockPagamentoRepository{}
	enviarUC := &EnviarPagamentoUseCaseImpl{repo: mockRepo}
	confirmarUC := &ConfirmarPagamentoUseCaseImpl{repo: mockRepo}

	pagamento, err := entities.PagamentoNew(1, 1, 100.00, "PENDENTE", "2024-05-15T10:00:00Z")
	if err != nil {
		t.Fatalf("erro ao criar pagamento: %v", err)
	}

	// Test sending the pagamento
	err = enviarUC.EnviarPagamento(context.Background(), pagamento)
	if err != nil {
		t.Fatalf("erro inesperado ao enviar pagamento: %v", err)
	}

	// Test confirming the pagamento
	pagamento.Status = "CONFIRMADO"
	err = confirmarUC.ConfirmarPagamento(context.Background(), pagamento)
	if err != nil {
		t.Fatalf("erro inesperado ao confirmar pagamento: %v", err)
	}

	// Check if the status was updated
	found := false
	for _, p := range mockRepo.pagamentos {
		if p.IdPagamento == pagamento.IdPagamento {
			found = true
			if p.Status != "CONFIRMADO" {
				t.Errorf("status do pagamento não foi atualizado: obtido %s, esperado CONFIRMADO", p.Status)
			}
		}
	}
	if !found {
		t.Errorf("pagamento não encontrado no repositório")
	}
}
