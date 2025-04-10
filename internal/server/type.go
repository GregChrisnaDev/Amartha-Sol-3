package server

import "github.com/GregChrisnaDev/Amartha-Sol-3/internal/handler"

type Handlers struct {
	UserHandler handler.UserHandler
	LoanHandler handler.LoanHandler
	LendHandler handler.LendHandler
}
