package main

import (
	"log"
	"net/http"
	"os"

	"github.com/GregChrisnaDev/Amartha-Sol-3/common/pdfgenerator"
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
	pdfGenerator := pdfgenerator.Init(os.Getenv("DEFAULT_TEMPLATE"), os.Getenv("DEFAULT_STORAGE"))
	transactionDB := postgres.NewDBTransaction(pgClient.DB)

	// init repository
	userRepo := repository.InitUserRepo(pgClient)
	loanRepo := repository.InitLoanRepo(pgClient)
	lendRepo := repository.InitLendRepo(pgClient)

	// init usecase
	userUC := usecase.InitUserUC(userRepo)
	loanUC := usecase.InitLoanUC(userRepo, loanRepo, lendRepo, storageClient, pdfGenerator)
	lendUC := usecase.InitLendUC(userRepo, loanRepo, lendRepo, transactionDB, storageClient, pdfGenerator)

	// init handler
	userHandler := handler.InitUserHandler(userUC)
	loanHandler := handler.InitLoanHandler(userUC, loanUC)
	lendHandler := handler.InitLendHandler(userUC, lendUC)

	r := server.RegisterRoute(server.Handlers{
		UserHandler: userHandler,
		LoanHandler: loanHandler,
		LendHandler: lendHandler,
	})

	// Start server
	log.Println("Server running on :9000")
	log.Fatal(http.ListenAndServe(":9000", r))
}
