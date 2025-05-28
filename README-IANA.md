# SWAGGO - Swagger for GO   

Instalar o CLI do SWAGGO

	go install github.com/swaggo/swag/cmd/swag@latest

Aplicar as alterações

	swag init 

	swag init --parseDependency --parseInternal
	
Rodar a aplicação com run go main.go
    
Acessar a doc na rota
        http://localhost:8080/docs/index.html

Para os atributos de TimeStamp e UltimaAtualizacao, utilizar o seguinte formato:

"2025-03-19T12:34:56Z"