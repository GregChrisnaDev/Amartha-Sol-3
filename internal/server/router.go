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
	loanRouter.HandleFunc("/propose", handlers.LoanHandler.ProposeLoanHandler).Methods("POST")
	loanRouter.HandleFunc("/get-all", handlers.LoanHandler.GetLoanByUIDHandler).Methods("GET")

	return r
}
