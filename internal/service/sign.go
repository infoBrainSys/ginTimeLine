package service

import "timeLineGin/internal/model"

type ISign interface {
	SignUp(*model.UserInput) error
	SignIn(*model.UserInput) error
}

var localSign ISign

func Sign() ISign {
	if localSign == nil {
		panic("localSign not register")
	}
	return localSign
}

func RegisterSign(s ISign) {
	localSign = s
}
