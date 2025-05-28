package route

import (
	"database/sql"
	"lanchonete/bootstrap"
	"lanchonete/internal/application/usecases"
	handler "lanchonete/internal/interfaces/http/handlers"
	repo "lanchonete/infra/database/repositories"

	"github.com/gin-gonic/gin"
)

// NewClienteRouter creates and configures all cliente-related routes
func NewClienteRouter(env *bootstrap.Env, db *sql.DB, router *gin.RouterGroup) {
	
	// Criar casos de uso
	clienteUseCase := usecases.NewClienteUseCase(repo.NewClienteMysqlRepository(db))

	cc := &handler.ClienteHandler{
		ClienteUseCase: clienteUseCase,
	}

	router.GET("/cliente/:CPF", cc.BuscarCliente)
	router.POST("/cliente", cc.CriarCliente)
}