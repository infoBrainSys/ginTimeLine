package sign

import (
	"timeLineGin/internal/model"
	"timeLineGin/internal/service"
	"timeLineGin/pkg/encrypt"
	sql "timeLineGin/pkg/mysql"
)

const (
	userExist = "用户已存在"
)

func init() {
	service.RegisterSign(New())
}

type SSign struct {
}

func New() *SSign {
	return new(SSign)
}

func (s *SSign) SignUp(u *model.UserInput) error {
	hashPass, err := encrypt.NewEncryptPassword().Encrypt(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashPass)

	input := model.UserInputCreate{
		UserInput: *u,
	}
	return sql.GetInstance().Save(&input).Error
}
