package server

import (
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/handler"
	"github.com/gorilla/mux"
)

func RegisterRoute(handlers Handlers) *mux.Router {
	r := mux.NewRouter()

	// health check
	r.HandleFunc("/ping", handler.PingHandler)

	// user route group
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/generate", handlers.UserHandler.UserGenerateHandler).Methods("POST")
	userRouter.HandleFunc("/get-all", handlers.UserHandler.GetAllUserHandler).Methods("GET")

	// loan route group
	loanRouter := r.PathPrefix("/loan").Subrouter()
	loanRouter.HandleFunc("/simulate", handlers.LoanHandler.SimulateLoanHandler).Methods("POST")
	loanRouter.HandleFunc("/propose", handlers.LoanHandler.ProposeLoanHandler).Methods("POST")
	loanRouter.HandleFunc("/get-all", handlers.LoanHandler.GetLoanByUIDHandler).Methods("GET")
	loanRouter.HandleFunc("/approve", handlers.LoanHandler.ApproveLoanHandler).Methods("POST")
	loanRouter.HandleFunc("/proof-pict", handlers.LoanHandler.GetProofPictureHandler).Methods("GET")
	loanRouter.HandleFunc("/list-approved-loan", handlers.LoanHandler.GetListApprovedLoanHandler).Methods("GET")
	loanRouter.HandleFunc("/disburse", handlers.LoanHandler.DisburseLoanHandler).Methods("POST")
	loanRouter.HandleFunc("/agreement-letter", handlers.LoanHandler.GetAgreementLetterHandler).Methods("GET")
	loanRouter.HandleFunc("/list-lender", handlers.LoanHandler.GetListLenderHandler).Methods("GET")

	// lend route group
	lendRoute := r.PathPrefix("/lend").Subrouter()
	lendRoute.HandleFunc("/simulate", handlers.LendHandler.SimulateHandler).Methods("POST")
	lendRoute.HandleFunc("/list-lend", handlers.LendHandler.GetListLenderHandler).Methods("GET")
	lendRoute.HandleFunc("/invest", handlers.LendHandler.InvestHandler).Methods("POST")
	lendRoute.HandleFunc("/agreement-letter", handlers.LendHandler.GetAgreementLetterHandler).Methods("GET")

	return r
}
