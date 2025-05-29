package handler

import (
	"fmt"
	_ "lanchonete/docs"
	"lanchonete/internal/application/presenters"
	"lanchonete/internal/domain/usecase"
	response "lanchonete/internal/interfaces/http/responses"
	"lanchonete/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AcompanhamentoHandler struct {
	AcompanhamentoUseCase        usecase.AcompanhamentoUseCase
	PedidoAtualizarStatusUseCase usecases.PedidoAtualizarStatusUseCase
}

func NewAcompanhamentoHandler(auc usecase.AcompanhamentoUseCase, p usecases.PedidoAtualizarStatusUseCase) *AcompanhamentoHandler {
	return &AcompanhamentoHandler{
		AcompanhamentoUseCase:        auc,
		PedidoAtualizarStatusUseCase: p,
	}
}

// CriarAcompanhamento godoc
// @Summary Cria um acompanhamento
// @Description Cria um acompanhamento
// @Tags acompanhamento
// @Router /acompanhamento [post]
// @Accept  json
// @Produce  json
// @Param acompanhamento body entities.AcompanhamentoPedido true "Acompanhamento"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
func (ah *AcompanhamentoHandler) CriarAcompanhamento(c *gin.Context) {
	id, err := ah.AcompanhamentoUseCase.CriarAcompanhamento(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse{
		Message: fmt.Sprintf("Acompanhamento criado com ID: %d", id),
	})
}

// AdicionarPedido godoc
// @Summary Adiciona um pedido ao acompanhamento
// @Description Adiciona um pedido existente ao acompanhamento de pedidos
// @Tags acompanhamento
// @Router /acompanhamento/{IDAcompanhamento}/{IDPedido} [post]
// @Accept json
// @Produce json
// @Param IDAcompanhamento path string true "ID do acompanhamento"
// @Param IDPedido path string true "ID do pedido"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse "Pedido ou acompanhamento não encontrado"
// @Failure 500 {object} response.ErrorResponse "Erro interno"
func (ah *AcompanhamentoHandler) AdicionarPedido(c *gin.Context) {
	idAcompanhamento := c.Param("IDAcompanhamento")
	idPedido := c.Param("IDPedido")

	idAcomp, err := strconv.Atoi(idAcompanhamento)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "IDAcompanhamento inválido"})
		return
	}
	idPed, err := strconv.Atoi(idPedido)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "IDPedido inválido"})
		return
	}

	err = ah.AcompanhamentoUseCase.AdicionarPedido(c, idAcomp, idPed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Pedido adicionado ao acompanhamento com sucesso",
	})
}

// BuscarAcompanhamento godoc
// @Summary Busca um acompanhamento
// @Description Busca um acompanhamento pelo ID
// @Tags acompanhamento
// @Router /acompanhamento/{ID} [get]
// @Accept json
// @Produce json
// @Param ID path string true "ID do acompanhamento"
// @Success 200 {object} presenters.AcompanhamentoDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse "Acompanhamento não encontrado"
func (ah *AcompanhamentoHandler) BuscarAcompanhamento(c *gin.Context) {
	idStr := c.Param("ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "ID inválido"})
		return
	}

	acompanhamento, err := ah.AcompanhamentoUseCase.BuscarAcompanhamento(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	dto := presenters.NewAcompanhamentoDTO(acompanhamento)
	c.JSON(http.StatusOK, dto)
}

// AtualizarStatusPedido godoc
// @Summary Atualiza o status de um pedido
// @Description Atualiza o status de um pedido no acompanhamento
// @Tags acompanhamento
// @Router /acompanhamento/{IDAcompanhamento}/pedido/{IDPedido}/status [put]
// @Accept json
// @Produce json
// @Param IDAcompanhamento path string true "ID do acompanhamento"
// @Param IDPedido path string true "ID do pedido"
// @Param status body StatusUpdateRequest true "Novo status do pedido"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse "Pedido ou acompanhamento não encontrado"
func (ah *AcompanhamentoHandler) AtualizarStatusPedido(c *gin.Context) {
	idPedidoStr := c.Param("IDPedido")
	fmt.Println("IDPedido: ", idPedidoStr)

	idPedido, err := strconv.Atoi(idPedidoStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "IDPedido inválido"})
		return
	}

	peds, err := ah.AcompanhamentoUseCase.BuscarPedidos(c, idPedido)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	fmt.Println("Pedidos encontrados: ", peds)

	c.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Status do pedido atualizado com sucesso",
	})
}

// BuscarPedidos godoc
// @Summary Busca os pedidos de um acompanhamento
// @Description Busca os pedidos associados a um acompanhamento
// @Tags acompanhamento
// @Router /acompanhamento/{ID}/pedidos [get]
// @Accept json
// @Produce json
// @Param ID path string true "ID do acompanhamento"
// @Success 200 {object} []entities.Pedido
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse "Acompanhamento não encontrado"
// @Failure 500 {object} response.ErrorResponse "Erro interno"
func (h *AcompanhamentoHandler) BuscarPedidos(ctx *gin.Context) {
	fmt.Println("Handler: ", ctx.Param("ID"))
	idAcompanhamentoStr := ctx.Param("ID")
	idAcompanhamento, err := strconv.Atoi(idAcompanhamentoStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "ID inválido"})
		return
	}

	pedidos, err := h.AcompanhamentoUseCase.BuscarPedidos(ctx, idAcompanhamento)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pedidos)
}

type StatusUpdateRequest struct {
	Status string `json:"status" example:"Em preparação"`
}
