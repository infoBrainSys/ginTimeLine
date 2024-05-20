package model

type UserInput struct {
	Nickname string `json:"nickname" form:"nickname" gorm:"type:string;not null;unique;comment:别名"`
	Passport string `json:"passport" form:"passport" binding:"required" gorm:"type:string;not null;unique;comment:账号"`
	Password string `json:"password" form:"password" binding:"required" gorm:"type:string;not null;comment:密码"`
	Email    string `json:"email" form:"email" gorm:"type:string;not null;unique;comment:邮箱"`
	Avatar   string `json:"avatar" form:"avatar" gorm:"type:string;not null;comment:头像"`
}

type UserInputCreate struct {
	Model
	UserInput
}

func (*UserInput) TableName() string {
	return "user"
}
