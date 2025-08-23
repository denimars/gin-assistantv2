package code

func BaseModel(projectName string) string {
	return `
	package model

import (
	"` + projectName + `/app/helper"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID string ` + "`gorm:\"type:varchar(50)\" json:\"id\"`" + `
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}

type DateTime struct {
	CreatedAt time.Time      ` + "`json:\"-\" gorm:\"type:datetime(0)\"`" + `
	UpdatedAt time.Time      ` + "`json:\"-\" gorm:\"type:datetime(0)\"`" + `
	DeletedAt gorm.DeletedAt ` + "`json:\"-\" gorm:\"type:datetime(0)\"`" + `
}

func (d *DateTime) BeforeSave(tx *gorm.DB) (err error) {
	d.CreatedAt = helper.UtcTime()
	d.UpdatedAt = helper.UtcTime()
	return
}

func (d *DateTime) BeforeUpdate(tx *gorm.DB) (err error) {
	d.UpdatedAt = helper.UtcTime()
	return
}

func (d *DateTime) BeforeDelete(tx *gorm.DB) (err error) {
	d.DeletedAt.Time = helper.UtcTime()
	return
}
	`
}
