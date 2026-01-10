package model

// RegisterReq user request to register handler.
type RegisterReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginReq user request to login handler.
type LoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
