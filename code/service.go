package code

import (
	"strings"

	"github.com/gin-assistantv2/helper"
)

func claerPackageName(serviceName string) string {
	return strings.NewReplacer(".", "", "/", "", "\\", "", "-", "").Replace(strings.ToLower(serviceName))
}

func Repository(name string) string {
	return `package ` + helper.GetServiceName(claerPackageName(name)) + `

import (
	"gorm.io/gorm"
)

type Repository interface {
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

`
}

func Service(name string) string {
	return `
package ` + helper.GetServiceName(claerPackageName(name)) + `

type Service interface {

}

type service_ struct {
	repository Repository
}

func NewService(repository Repository) *service_ {
	return &service_{repository}
}

`
}

func Router(name string) string {
	packageName := helper.GetServiceName(claerPackageName(name))
	return `
package ` + packageName + `

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(g *gin.RouterGroup, db *gorm.DB) {
	//repository := NewRepository(db)
	//service_ := NewService(repository)
	//route := g.Group("/` + packageName + `")

	//route.GET("",)

}

`
}
