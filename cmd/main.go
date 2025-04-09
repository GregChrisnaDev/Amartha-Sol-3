package main

import (
	"log"
	"net/http"
	"os"

	"github.com/GregChrisnaDev/Amartha-Sol-3/common/postgres"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/handler"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/repository"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/server"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/storage"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/usecase"
)

func main() {

	// init client
	pgClient := postgres.Init()
	storageClient := storage.Init(os.Getenv("DEFAULT_STORAGE"))

	// init repository
	userRepo := repository.InitUserRepo(pgClient)
	loanRepo := repository.InitLoanRepo(pgClient)

	// init usecase
	userUC := usecase.InitUserUC(userRepo)
	loanUC := usecase.InitLoanUC(loanRepo, storageClient)

	// init handler
	userHandler := handler.InitUserHandler(userUC)
	loanHandler := handler.InitLoanHandler(userUC, loanUC)

	r := server.RegisterRoute(server.Handlers{
		UserHandler: userHandler,
		LoanHandler: loanHandler,
	})

	// Start server
	log.Println("Server running on :9000")
	log.Fatal(http.ListenAndServe(":9000", r))
}
