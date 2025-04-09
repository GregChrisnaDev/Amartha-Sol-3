package model

type User struct {
	ID           uint64
	Name         string
	Address      string
	Role         int
	PasswordHash string
}

// User Role
const (
	Customer = iota + 1
	Employee
)

var RoleMap = map[int]string{
	Customer: "Customer",
	Employee: "Employee",
}
