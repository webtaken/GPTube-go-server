package models

type RegisterReq struct {
	Email string `json:"email,omitempty"`
}

type RegisterResp struct {
	Err string `json:"error,omitempty"`
}
