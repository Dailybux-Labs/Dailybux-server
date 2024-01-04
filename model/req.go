package model

type LoginReq struct {
	NickName   string `json:"nickName"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Wallet     string `json:"wallet"`
	Google     string `json:"google"`
	MobileType int    `json:"mobileType"`
}

type UserInfoReq struct {
	Id uint `json:"id"`
}
