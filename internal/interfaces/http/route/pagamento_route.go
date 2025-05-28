package route

import (
	"database/sql"
	"lanchonete/bootstrap"
	"lanchonete/internal/application/usecases"
	repo "lanchonete/infra/database/repositories"
	handler "lanchonete/internal/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
)

// NewPagamentoRouter creates and configures all pagamento-related routes
func NewPagamentoRouter(env *bootstrap.Env, db *sql.DB, router *gin.RouterGroup) {
	pagamentoRepo := repo.NewPagamentoMysqlRepository(db)

	pc := &handler.PagamentoHandler{
		EnviarPagamentoUseCase: usecases.NewEnviarPagamentoUseCase(pagamentoRepo),
		ConfirmarPagamentoUseCase: usecases.NewConfirmarPagamentoUseCase(pagamentoRepo),
	}

	// Register payment routes
	router.POST("/pagamento", pc.EnviarPagamento)
	router.POST("/pagamento/confirmar", pc.ConfirmarPagamento)
}