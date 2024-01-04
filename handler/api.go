package handler

import (
	"dailybux/db/repo"
	"dailybux/model"
)

func Login(req *model.LoginReq) (res any, err error) {
	var userInfo model.UserInfo
	userInfo.Google = req.Google
	userInfo.Email = req.Email
	userInfo.Wallet = req.Wallet
	userInfo.MobileType = req.MobileType
	userInfo.NickName = req.NickName
	userInfo.Phone = req.Phone
	query, err := repo.UserQuery(&userInfo)
	if err != nil {
		return nil, err
	}
	if query.Id == 0 {
		fullUserInfo, err := repo.UserInsert(&userInfo)
		if err != nil {
			return nil, err
		}
		return fullUserInfo, nil
	}
	return query, nil
}

func UserInfo(req *model.UserInfoReq) (res any, err error) {
	fullUserInfo, err := repo.UserQueryById(req.Id)
	if err != nil {
		return nil, err
	}
	return fullUserInfo, nil
}

func DailyCheckIn(req *model.UserInfoReq) (res any, err error) {
	return 0, nil
}

func Crunch() (res any, err error) {
	var crunchInfo model.CrunchInfo
	crunchInfo.DailyRewards = "200"
	crunchInfo.ExchangeRate = "0.0001"
	return crunchInfo, nil
}

func Peanut() (res any, err error) {
	return nil, nil
}
