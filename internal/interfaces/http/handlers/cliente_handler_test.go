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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClienteUseCase is a mock implementation of usecases.ClienteUseCase
type MockClienteUseCase struct {
	mock.Mock
}

func (m *MockClienteUseCase) CriarCliente(ctx context.Context, cliente *entities.Cliente) error {
	args := m.Called(ctx, cliente)
	return args.Error(0)
}

func (m *MockClienteUseCase) BuscarCliente(ctx context.Context, cpf string) (*entities.Cliente, error) {
	args := m.Called(ctx, cpf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Cliente), args.Error(1)
}

func setupClienteTestRouter() (*gin.Engine, *MockClienteUseCase) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockCliente := new(MockClienteUseCase)

	handler := &ClienteHandler{
		ClienteUseCase: mockCliente,
	}

	router.POST("/cliente", handler.CriarCliente)
	router.GET("/cliente/:CPF", handler.BuscarCliente)

	return router, mockCliente
}

func TestCriarCliente(t *testing.T) {
	router, mockCliente := setupClienteTestRouter()

	tests := []struct {
		name           string
		cliente        entities.Cliente
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful creation",
			cliente: entities.Cliente{
				Nome:  "Test Client",
				Email: "test@example.com",
				CPF:   "12345678900",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "error creating cliente",
			cliente: entities.Cliente{
				Nome:  "Test Client",
				Email: "test@example.com",
				CPF:   "12345678900",
			},
			mockError:      errors.New("error creating cliente"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock for each test
			mockCliente.ExpectedCalls = nil
			mockCliente.Calls = nil

			// Set up mock expectation
			mockCliente.On("CriarCliente", mock.Anything, mock.MatchedBy(func(cliente *entities.Cliente) bool {
				return cliente.Nome == tt.cliente.Nome &&
					cliente.Email == tt.cliente.Email &&
					cliente.CPF == tt.cliente.CPF
			})).Return(tt.mockError)

			payloadBytes, _ := json.Marshal(tt.cliente)
			req := httptest.NewRequest(http.MethodPost, "/cliente", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockCliente.AssertExpectations(t)
		})
	}
}

func TestBuscarCliente(t *testing.T) {
	router, mockCliente := setupClienteTestRouter()

	tests := []struct {
		name           string
		cpf            string
		mockCliente    *entities.Cliente
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful fetch",
			cpf:  "12345678900",
			mockCliente: &entities.Cliente{
				Nome:  "Test Client",
				Email: "test@example.com",
				CPF:   "12345678900",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "cliente not found",
			cpf:            "12345678900",
			mockCliente:    nil,
			mockError:      errors.New("cliente not found"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock for each test
			mockCliente.ExpectedCalls = nil
			mockCliente.Calls = nil

			// Set up mock expectation
			mockCliente.On("BuscarCliente", mock.Anything, tt.cpf).Return(tt.mockCliente, tt.mockError)

			req := httptest.NewRequest(http.MethodGet, "/cliente/"+tt.cpf, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockCliente.AssertExpectations(t)
		})
	}
}
