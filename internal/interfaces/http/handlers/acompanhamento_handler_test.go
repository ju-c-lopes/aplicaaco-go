package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"lanchonete/internal/domain/entities"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock UseCases ---

type MockAcompanhamentoUseCase struct{ mock.Mock }

// Implements: AdicionarPedido(c context.Context, idAcompanhamento int, pedido *entities.Pedido) error
func (m *MockAcompanhamentoUseCase) AdicionarPedido(c context.Context, idAcompanhamento int, pedidoID int) error {
	args := m.Called(c, idAcompanhamento, pedidoID)
	return args.Error(0)
}

// Implements: BuscarAcompanhamento(c context.Context, id int) (*entities.AcompanhamentoPedido, error)
func (m *MockAcompanhamentoUseCase) BuscarAcompanhamento(c context.Context, id int) (*entities.AcompanhamentoPedido, error) {
	args := m.Called(c, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AcompanhamentoPedido), args.Error(1)
}

// Implements: AtualizarStatusPedido(c context.Context, idAcompanhamento int, status entities.StatusPedido) error
func (m *MockAcompanhamentoUseCase) AtualizarStatusPedido(c context.Context, idAcompanhamento int, status entities.StatusPedido) error {
	args := m.Called(c, idAcompanhamento, status)
	return args.Error(0)
}

// Implements: CriarAcompanhamento(c context.Context, acompanhamento *entities.AcompanhamentoPedido) error
func (m *MockAcompanhamentoUseCase) CriarAcompanhamento(c context.Context) (int, error) {
	return 0, nil
}

// --- Setup Handler ---

func setupAcompanhamentoHandlerWithMocks() (*AcompanhamentoHandler, *MockAcompanhamentoUseCase, *MockPedidoAtualizarStatusUseCase) {
	mockAcompanhamento := new(MockAcompanhamentoUseCase)
	mockPedidoAtualizar := new(MockPedidoAtualizarStatusUseCase)

	handler := &AcompanhamentoHandler{
		AcompanhamentoUseCase:        mockAcompanhamento,
		PedidoAtualizarStatusUseCase: mockPedidoAtualizar,
	}
	return handler, mockAcompanhamento, mockPedidoAtualizar
}

func TestAcompanhamentoHandler_AdicionarPedido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockAcompanhamento, _ := setupAcompanhamentoHandlerWithMocks()

	mockAcompanhamento.On("AdicionarPedido", mock.Anything, "acomp1", mock.AnythingOfType("*entities.Pedido")).Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "/acompanhamento/acomp1/ped1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "IDAcompanhamento", Value: "acomp1"},
		{Key: "IDPedido", Value: "ped1"},
	}
	c.Request = req

	handler.AdicionarPedido(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Pedido adicionado ao acompanhamento com sucesso")
}

func TestAcompanhamentoHandler_BuscarAcompanhamento(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockAcompanhamento, _ := setupAcompanhamentoHandlerWithMocks()

	acomp := &entities.AcompanhamentoPedido{ID: 1}
	mockAcompanhamento.On("BuscarAcompanhamento", mock.Anything, 1).Return(acomp, nil)

	req, _ := http.NewRequest(http.MethodGet, "/acompanhamento/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "ID", Value: "1"}}
	c.Request = req

	handler.BuscarAcompanhamento(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAcompanhamentoHandler_AtualizarStatusPedido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockAcompanhamento, _ := setupAcompanhamentoHandlerWithMocks()

	statusReq := StatusUpdateRequest{Status: "Finalizado"}
	mockAcompanhamento.On("AtualizarStatusPedido", mock.Anything, "acomp1", entities.StatusPedido("Finalizado")).Return(nil)

	body, _ := json.Marshal(statusReq)
	req, _ := http.NewRequest(http.MethodPut, "/acompanhamento/acomp1/pedido/ped1/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "IDAcompanhamento", Value: "acomp1"},
		{Key: "IDPedido", Value: "ped1"},
	}
	c.Request = req

	handler.AtualizarStatusPedido(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Status do pedido atualizado com sucesso")
}
