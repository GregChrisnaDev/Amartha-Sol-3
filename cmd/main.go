package main

import (
	"log"
	"net/http"

	"github.com/GregChrisnaDev/Amartha-Sol-3/common/postgres"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/handler"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/repository"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/server"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/usecase"
)

func main() {

	// init client
	pgClient := postgres.Init()

	// init repository
	userRepo := repository.InitUserRepo(pgClient)

	// init usecase
	userUC := usecase.InitUserUC(userRepo)

	// init handler
	userHandler := handler.InitUserHandler(userUC)

	r := server.RegisterRoute(server.Handlers{
		UserHandler: userHandler,
	})

	// Start server
	log.Println("Server running on :9000")
	log.Fatal(http.ListenAndServe(":9000", r))
}
