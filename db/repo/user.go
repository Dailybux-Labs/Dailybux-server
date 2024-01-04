package repo

import (
	"errors"

	"gorm.io/gorm"

	"dailybux/config"
	"dailybux/model"
)

func UserInsert(userInfo *model.UserInfo) (id uint, err error) {
	if err := config.Db0.Table("user_info").Create(&userInfo).Error; err != nil {
		return 0, err
	}
	return userInfo.Id, nil
}

func UserQuery(userInfo *model.UserInfo) (fullUserInfo *model.UserInfo, err error) {
	if err := config.Db0.Table("user_info").Where("phone = ? or email = ? or wallet = ? or google = ?", userInfo.Phone, userInfo.Email, userInfo.Wallet, userInfo.Google).Find(&fullUserInfo).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return fullUserInfo, nil
}

func UserQueryById(id uint) (fullUserInfo *model.UserInfo, err error) {
	if err := config.Db0.Table("user_info").Where("id = ?", id).Find(&fullUserInfo).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return fullUserInfo, nil
}
