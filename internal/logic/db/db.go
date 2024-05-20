package db

import (
	"errors"
	"github.com/gogf/gf/v2/util/gconv"
	"gorm.io/gorm"
	"timeLineGin/internal/model"
	"timeLineGin/pkg/encrypt"
	"timeLineGin/pkg/mysql"
)

type DB struct {
	D     *gorm.DB
	value string
	m     any
}

func NewDB() *DB {
	return &DB{
		D: mysql.GetInstance(),
	}
}

func (d *DB) Get(input *model.UserInput) *DB {
	var u model.UserInputCreate
	if err := d.D.Model(&model.UserInput{}).
		Where("passport = ?", input.Passport).
		Scan(&u).
		Error; err != nil {
		return d
	}
	d.m = &u
	return d
}

// UserExist 查询用户
func (d *DB) UserExist(passport string) error {
	err := d.D.Model(&model.UserInput{}).
		First("passport = ?", passport).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Exist 判断是否存在
func (d *DB) Exist() bool {
	if d.value == "" {
		return false
	}
	return true
}

// Column 比对某个字段，如果存在则赋值给 d.value
func (d *DB) Column(column string) *DB { // TODO 函数未实现

	d.value = ""
	m := d.m.(*model.UserInputCreate)
	mp := gconv.Map(&m)
	for i, v := range mp {
		if mp[i] == column {
			d.value = v.(string)
			return d
		}
	}
	return d
}

func (d *DB) ComparePassword(input string) error {
	return func() error {
		d.value = ""
		m := d.m.(*model.UserInputCreate)
		err := encrypt.NewEncryptPassword().ComparePassword([]byte(m.Password), input)
		return err
	}()
}

func (d *DB) Value() string {
	return d.value
}
