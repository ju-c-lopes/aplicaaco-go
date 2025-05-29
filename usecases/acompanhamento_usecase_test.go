package usecases

import (
	"context"
	"testing"
	"time"

	"lanchonete/internal/domain/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock Implementation ---

type MockAcompanhamentoUseCase struct {
	mock.Mock
}

func (m *MockAcompanhamentoUseCase) CriarAcompanhamento(c context.Context, acompanhamento *entities.AcompanhamentoPedido) error {
	args := m.Called(c, acompanhamento)
	return args.Error(0)
}
func (m *MockAcompanhamentoUseCase) BuscarPedidos(c context.Context, ID string) (entities.Pedido, error) {
	args := m.Called(c, ID)
	return args.Get(0).(entities.Pedido), args.Error(1)
}
func (m *MockAcompanhamentoUseCase) AdicionarPedido(c context.Context, acompanhamento *entities.AcompanhamentoPedido, pedido *entities.Pedido) error {
	args := m.Called(c, acompanhamento, pedido)
	return args.Error(0)
}
func (m *MockAcompanhamentoUseCase) BuscarAcompanhamento(c context.Context, ID string) (*entities.AcompanhamentoPedido, error) {
	args := m.Called(c, ID)
	return args.Get(0).(*entities.AcompanhamentoPedido), args.Error(1)
}
func (m *MockAcompanhamentoUseCase) AtualizarStatusPedido(c context.Context, acompanhamentoID string, identificacao string, novoStatus entities.StatusPedido) error {
	args := m.Called(c, acompanhamentoID, identificacao, novoStatus)
	return args.Error(0)
}

// --- Test ---

func TestCriarAcompanhamento_Completo(t *testing.T) {
	mockUC := new(MockAcompanhamentoUseCase)
	ctx := context.Background()

	acomp := &entities.AcompanhamentoPedido{
		ID:                1,
		Pedidos:           []entities.Pedido{},
		TempoEstimado:     "00:15:00", // 15 minutes
		UltimaAtualizacao: time.Now(),
	}

	mockUC.On("CriarAcompanhamento", ctx, acomp).Return(nil)

	err := mockUC.CriarAcompanhamento(ctx, acomp)
	assert.NoError(t, err)
	assert.Equal(t, "acomp-1", acomp.ID)
	assert.Equal(t, 20*time.Minute, acomp.TempoEstimado)
	assert.NotEmpty(t, acomp.UltimaAtualizacao)
}
