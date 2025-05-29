package route

import (
	"database/sql"
	"lanchonete/bootstrap"
	appusecases "lanchonete/internal/application/usecases"
	handler "lanchonete/internal/interfaces/http/handlers"
	repo "lanchonete/infra/database/repositories"
	"lanchonete/usecases"

	"fmt"

	"github.com/gin-gonic/gin"
)

// NewAcompanhamentoRouter creates and configures all acompanhamento-related routes
func NewAcompanhamentoRouter(env *bootstrap.Env, db *sql.DB, router *gin.RouterGroup) {

	// Criar casos de uso
	acompanhamentoUseCase := appusecases.NewAcompanhamentoUseCase(repo.NewAcompanhamentoMySQLRepository(db))
	pedidoAtualizarStatusUseCase := usecases.NewPedidoAtualizarStatusUseCase(repo.NewPedidoMysqlRepository(db))

	auc := &handler.AcompanhamentoHandler{
		AcompanhamentoUseCase:        acompanhamentoUseCase,
		PedidoAtualizarStatusUseCase: pedidoAtualizarStatusUseCase,
	}

	fmt.Printf("Registrando rotas do acompanhamento\n")

	router.POST("/acompanhamento", auc.CriarAcompanhamento)
	router.GET("/acompanhamento/:ID", auc.BuscarAcompanhamento)
	router.POST("/acompanhamento/:IDAcompanhamento/:IDPedido", auc.AdicionarPedido)
	router.PUT("acompanhamento/:IDAcompanhamento/:IDPedido/:status", auc.AtualizarStatusPedido)
	router.GET("/acompanhamento/:ID/pedidos", auc.BuscarPedidos)
}
