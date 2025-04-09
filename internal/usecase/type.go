package usecase

type UserGenerateReq struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type UserResp struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Role     string `json:"role"`
	Password string `json:"password,omitempty"` // this only use for testing purpose to simplify when needed
}

type ValidateUserReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
