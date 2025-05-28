package server

import (
	"lanchonete/bootstrap"
	"lanchonete/internal/interfaces/http/handlers"
	"lanchonete/usecases"
	usecs "lanchonete/internal/application/usecases"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

type Server struct {
	app    *bootstrap.App
	router *gin.Engine
}

func NewServer(app *bootstrap.App) *Server {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	return &Server{
		app:    app,
		router: router,
	}
}

func (s *Server) SetupRoutes() {
	api := s.router.Group("")

	// Cliente
	clienteUseCase := usecs.NewClienteUseCase(s.app.ClienteRepository)
	clienteHandler := &handler.ClienteHandler{
		ClienteUseCase: clienteUseCase,
	}
	api.POST("/cliente", clienteHandler.CriarCliente)
	api.GET("/cliente/:CPF", clienteHandler.BuscarCliente)

	// Produto
	produtoRepo := s.app.ProdutoRepository
	produtoIncluir := usecases.NewProdutoIncluirUseCase(produtoRepo)
	produtoEditar := usecases.NewProdutoEditarUseCase(produtoRepo)
	produtoRemover := usecases.NewProdutoRemoverUseCase(produtoRepo)
	produtoBuscar := usecases.NewProdutoBuscaPorIdUseCase(produtoRepo)
	produtoListarTodos := usecases.NewProdutoListarTodosUseCase(produtoRepo)
	produtoListarPorCategoria := usecases.NewProdutoListarPorCategoriaUseCase(produtoRepo)

	produtoHandler := handler.NewProdutoHandler(
		produtoIncluir,
		produtoBuscar,
		produtoListarTodos,
		produtoEditar,
		produtoRemover,
		produtoListarPorCategoria,
	)
	api.POST("/produtos", produtoHandler.ProdutoIncluir)
	api.GET("/produtos/:id", produtoHandler.ProdutoBuscarPorId)
	api.GET("/produtos", produtoHandler.ProdutoListarTodos)
	api.PUT("/produtos/editar", produtoHandler.ProdutoEditar)
	api.DELETE("/produtos/:id", produtoHandler.ProdutoRemover)
	api.GET("/produtos/categoria/:categoria", produtoHandler.ProdutoListarPorCategoria)

	// Pedido
	pedidoRepo := s.app.PedidoRepository
	pedidoIncluir := usecases.NewPedidoIncluirUseCase(pedidoRepo)
	pedidoBuscar := usecases.NewPedidoBuscarPorIdUseCase(pedidoRepo)
	pedidoAtualizar := usecases.NewPedidoAtualizarStatusUseCase(pedidoRepo)
	pedidoListarTodos := usecases.NewPedidoListarTodosUseCase(pedidoRepo)
	produtoBuscaPorId := usecases.NewProdutoBuscaPorIdUseCase(produtoRepo)

	pedidoHandler := handler.NewPedidoHandler(
		pedidoIncluir,
		pedidoBuscar,
		pedidoAtualizar,
		produtoBuscaPorId,
		pedidoListarTodos,
	)
	api.POST("/pedidos", pedidoHandler.CriarPedido)
	api.GET("/pedidos/:nroPedido", pedidoHandler.BuscarPedido)
	api.PUT("/pedidos/:nroPedido/status/:status", pedidoHandler.AtualizarStatusPedido)
	api.POST("/pedidos/listartodos", pedidoHandler.ListarTodosOsPedidos)

	// Acompanhamento
	acompUseCase := usecs.NewAcompanhamentoUseCase(s.app.AcompanhamentoRepository)
	acompHandler := handler.NewAcompanhamentoHandler(acompUseCase, pedidoAtualizar)

	api.POST("/acompanhamento", acompHandler.CriarAcompanhamento)
	api.POST("/acompanhamento/:IDAcompanhamento/:IDPedido", acompHandler.AdicionarPedido)
	api.GET("/acompanhamento/:ID", acompHandler.BuscarAcompanhamento)
	api.PUT("/acompanhamento/:IDAcompanhamento/pedido/:IDPedido/status", acompHandler.AtualizarStatusPedido)

	// Pagamento — se implementado
	if s.app.PagamentoRepository != nil {
		// Exemplo:
		// pagamentoUseCase := usecases.NewPagamentoUseCase(s.app.PagamentoRepository)
		// pagamentoHandler := handler.NewPagamentoHandler(pagamentoUseCase)
		// api.POST("/pagamento", pagamentoHandler.CriarPagamento)
	}

	// Swagger Docs
	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (s *Server) Start() error {
	return s.router.Run(s.app.Env.ServerAddress)
}
