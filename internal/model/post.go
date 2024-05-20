package model

type PostInput struct {
	Title    string `json:"title" form:"title" binding:"required" gorm:"type:string;not null;comment:标题"`
	Content  string `json:"content" form:"content" binding:"required" gorm:"type:string;not null;comment:内容"`
	Category string `json:"category" form:"category" binding:"required" gorm:"type:string;not null;comment:分类"`
}

type PostInputCreate struct {
	*Model
	*PostInput
	UserID uint `json:"user_id" form:"user_id" binding:"required" gorm:"type:int;not null;comment:用户ID"`
}

func (p *PostInputCreate) TableName() string {
	return "post"
}
