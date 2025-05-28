package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"lanchonete/internal/domain/entities"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEnviarPagamentoUseCase is a mock implementation of usecases.EnviarPagamentoUseCase
type MockEnviarPagamentoUseCase struct {
	mock.Mock
}

func (m *MockEnviarPagamentoUseCase) EnviarPagamento(ctx context.Context, pagamento *entities.Pagamento) error {
	args := m.Called(ctx, pagamento)
	return args.Error(0)
}

// MockConfirmarPagamentoUseCase is a mock implementation of usecases.ConfirmarPagamentoUseCase
type MockConfirmarPagamentoUseCase struct {
	mock.Mock
}

func (m *MockConfirmarPagamentoUseCase) ConfirmarPagamento(ctx context.Context, pagamento *entities.Pagamento) error {
	args := m.Called(ctx, pagamento)
	return args.Error(0)
}

func setupPagamentoTestRouter() (*gin.Engine, *MockEnviarPagamentoUseCase, *MockConfirmarPagamentoUseCase) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockEnviarPagamento := new(MockEnviarPagamentoUseCase)
	mockConfirmarPagamento := new(MockConfirmarPagamentoUseCase)

	handler := &PagamentoHandler{
		EnviarPagamentoUseCase:    mockEnviarPagamento,
		ConfirmarPagamentoUseCase: mockConfirmarPagamento,
	}

	router.POST("/pagamento", handler.EnviarPagamento)
	router.POST("/pagamento/confirmar", handler.ConfirmarPagamento)

	return router, mockEnviarPagamento, mockConfirmarPagamento
}

func TestEnviarPagamento(t *testing.T) {
	router, mockEnviarPagamento, _ := setupPagamentoTestRouter()

	tests := []struct {
		name           string
		pagamento      entities.Pagamento
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful payment sending",
			pagamento: entities.Pagamento{
				IdPagamento:  123,
				Status:      "pending",
				Valor:       100.00,
				DataCriacao: time.Now().Format("2006-01-02 15:04:05"),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "error sending payment",
			pagamento: entities.Pagamento{
				IdPagamento:  123,
				Status:      "pending",
				Valor:       100.00,
				DataCriacao: time.Now().Format("2006-01-02 15:04:05"),
			},
			mockError:      errors.New("error sending payment"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock for each test
			mockEnviarPagamento.ExpectedCalls = nil
			mockEnviarPagamento.Calls = nil

			// Set up mock expectation
			mockEnviarPagamento.On("EnviarPagamento", mock.Anything, mock.MatchedBy(func(pagamento *entities.Pagamento) bool {
				return pagamento.IdPagamento == tt.pagamento.IdPagamento &&
					pagamento.Status == tt.pagamento.Status &&
					pagamento.Valor == tt.pagamento.Valor &&
					pagamento.DataCriacao == tt.pagamento.DataCriacao
			})).Return(tt.mockError)

			payloadBytes, _ := json.Marshal(tt.pagamento)
			req := httptest.NewRequest(http.MethodPost, "/pagamento", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockEnviarPagamento.AssertExpectations(t)
		})
	}
}

func TestConfirmarPagamento(t *testing.T) {
	router, _, mockConfirmarPagamento := setupPagamentoTestRouter()

	tests := []struct {
		name           string
		pagamento      entities.Pagamento
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful payment confirmation",
			pagamento: entities.Pagamento{
				IdPagamento:  123,
				Status:      "confirmed",
				Valor:       100.00,
				DataCriacao: time.Now().Format("2006-01-02 15:04:05"),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "error confirming payment",
			pagamento: entities.Pagamento{
				IdPagamento:  123,
				Status:      "confirmed",
				Valor:       100.00,
				DataCriacao: time.Now().Format("2006-01-02 15:04:05"),
			},
			mockError:      errors.New("error confirming payment"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock for each test
			mockConfirmarPagamento.ExpectedCalls = nil
			mockConfirmarPagamento.Calls = nil

			// Set up mock expectation
			mockConfirmarPagamento.On("ConfirmarPagamento", mock.Anything, mock.MatchedBy(func(pagamento *entities.Pagamento) bool {
				return pagamento.IdPagamento == tt.pagamento.IdPagamento &&
					pagamento.Status == tt.pagamento.Status &&
					pagamento.Valor == tt.pagamento.Valor &&
					pagamento.DataCriacao == tt.pagamento.DataCriacao
			})).Return(tt.mockError)

			payloadBytes, _ := json.Marshal(tt.pagamento)
			req := httptest.NewRequest(http.MethodPost, "/pagamento/confirmar", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockConfirmarPagamento.AssertExpectations(t)
		})
	}
} 