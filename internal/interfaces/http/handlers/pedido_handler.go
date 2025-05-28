package handler

import (
	"encoding/json"
	"fmt"
	_ "lanchonete/docs"
	"lanchonete/internal/domain/entities"
	response "lanchonete/internal/interfaces/http/responses"
	"lanchonete/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PedidoHandler struct {
	PedidoIncluirUseCase         usecases.PedidoIncluirUseCase
	PedidoBuscarPorIdUseCase     usecases.PedidoBuscarPorIdUseCase
	PedidoAtualizarStatusUseCase usecases.PedidoAtualizarStatusUseCase
	ProdutoBuscarPorIdUseCase    usecases.ProdutoBuscaPorIdUseCase
	PedidoListarTodosUseCase     usecases.PedidoListarTodosUseCase
}

func NewPedidoHandler(pedidoIncluirUseCase usecases.PedidoIncluirUseCase,
	pedidoBuscarPorIdUseCase usecases.PedidoBuscarPorIdUseCase,
	pedidoAtualizarStatusUsecase usecases.PedidoAtualizarStatusUseCase,
	produtoBuscarPorIdUseCase usecases.ProdutoBuscaPorIdUseCase,
	pedidoListarTodosUseCase usecases.PedidoListarTodosUseCase) *PedidoHandler {
	return &PedidoHandler{
		PedidoIncluirUseCase:         pedidoIncluirUseCase,
		PedidoBuscarPorIdUseCase:     pedidoBuscarPorIdUseCase,
		PedidoAtualizarStatusUseCase: pedidoAtualizarStatusUsecase,
		ProdutoBuscarPorIdUseCase:    produtoBuscarPorIdUseCase,
		PedidoListarTodosUseCase:     pedidoListarTodosUseCase,
	}
}

// CriarPedido godoc
// @Summary Cria um pedido
// @Description Cria um pedido
// @Tags pedido
// @Router /pedidos [post]
// @Accept  json
// @Produce  json
// @Param pedido body entities.Pedido true "Pedido"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
func (h *PedidoHandler) CriarPedido(r *gin.Context) {
	var pedido entities.Pedido
	fmt.Println("Handler Criando pedido", pedido)
	err := json.NewDecoder(r.Request.Body).Decode(&pedido)
	fmt.Println("Handler Criando Depois pedido", pedido)
	if err != nil {
		r.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}
	
	// Substituir o array de produtos com os dados completos do banco
	produtosCompletos := []entities.Produto{}

	for _, produto := range pedido.Produtos {
		fmt.Println("Handler Buscando produto:", produto.Nome)
		pBanco, err := h.ProdutoBuscarPorIdUseCase.Run(r, produto.ID)
		if err != nil {
			r.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Produto não Cadastrado!"})
			return
		}
		produtosCompletos = append(produtosCompletos, *pBanco)
	}

	// Chamar PedidoNew com os produtos completos
	ped, err := h.PedidoIncluirUseCase.Run(r, pedido.ClienteCPF, produtosCompletos)
	if err != nil {
		r.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	r.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Pedido criado com sucesso" + strconv.Itoa(ped.ID),
	})
}

// BuscarPedido godoc
// @Summary Busca um pedido
// @Description Busca um pedido
// @Tags pedido
// @Router /pedidos/{ID} [get]
// @Accept  json
// @Produce  json
// @Param ID path string true "Número do pedido"
// @Success 200 {object} entities.Pedido
// @Failure 400 {object} response.ErrorResponse
func (h *PedidoHandler) BuscarPedido(r *gin.Context) {
	nroPedido := r.Param("nroPedido")
	id, err := strconv.Atoi(nroPedido)
	if err != nil {
		r.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Número do pedido inválido"})
		return
	}
	pedido, err := h.PedidoBuscarPorIdUseCase.Run(r, id)
	if err != nil {
		r.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	r.JSON(http.StatusOK, pedido)

}

// BuscarPedido godoc
// @Summary Atualiza um pedido a partir de sua Identificação
// @Description Atualizar um pedido
// @Tags pedido
// @Router /pedidos/{nroPedido}/status/{status} [put]
// @Accept  json
// @Produce  json
// @Param nroPedido path string true "Número do pedido"
// @Param status path string true "Novo Status do pedido"
// @Success 200 {object} entities.Pedido
// @Failure 400 {object} response.ErrorResponse
func (h *PedidoHandler) AtualizarStatusPedido(r *gin.Context) {
	nroPedido := r.Param("nroPedido")
	id, err := strconv.Atoi(nroPedido)
	status := r.Param("status")
	err = h.PedidoAtualizarStatusUseCase.Run(r, id, status)
	if err != nil {
		r.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	r.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Status do pedido atualizado com sucesso",
	})
}

// ProdutoListarTodos godoc
// @Summary Lista todos os pedidos no banco
// @Description Lista todos os pedidos presentes no banco
// @Tags pedido
// @Router /pedidos/listartodos [POST]
// @Accept  json
// @Produce  json
// @Success 200 {object} []entities.Pedido
// @Failure 400 {object} response.ErrorResponse
func (h *PedidoHandler) ListarTodosOsPedidos(r *gin.Context) {
	pedidos, err := h.PedidoListarTodosUseCase.Run(r)
	if err != nil {
		r.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	r.JSON(http.StatusOK, pedidos)
}
