package sign

import (
	"errors"
	"timeLineGin/internal/logic/db"
	"timeLineGin/internal/model"
)

const (
	userNotExist      = "user not exist"
	userPasswordError = "user password error"
)

func (s *SSign) SignIn(input *model.UserInput) error {
	d := db.NewDB()
	// 查找是否存在用户
	err := d.UserExist(input.Passport)
	if err != nil {
		return errors.New(userNotExist)
	}

	// 如果存在则校验密码是否正确
	if err := d.Get(input).ComparePassword(input.Password); err != nil {
		return errors.New(userPasswordError)
	}
	return nil
}
