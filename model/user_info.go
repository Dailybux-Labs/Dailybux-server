package model

type UserInfo struct {
	Id         uint   `json:"id"`
	NickName   string `json:"nickName"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Wallet     string `json:"wallet"`
	Google     string `json:"google"`
	MobileType int    `json:"mobileType"`
}
