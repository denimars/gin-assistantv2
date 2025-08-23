package code

func ModelBlackListToken() string {
	return `
package model

type BlackListToken struct {
	BaseModel
	Token string ` + "`gorm:\"type:varchar(255);not null\"`" + `
	DateTime
}

func (BlackListToken) TableName() string {
	return "black_list_tokens"
}
`
}

func RepoBlackListToken(projectName string) string {
	return `
package blacklisttoken

import (
	"` + projectName + `/app/model"

	"gorm.io/gorm"
)

type Repository interface {
	Create(token string) error
	FindByToken(token string) model.BlackListToken
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(token string) error {
	err := r.db.Debug().Create(&model.BlackListToken{
		Token: token,
	}).Error
	return err
}

func (r *repository) FindByToken(token string) model.BlackListToken {
	var blacklistToken model.BlackListToken
	r.db.Debug().Where("token = ?", token).Find(&blacklistToken)
	return blacklistToken
}

`
}
