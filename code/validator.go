package code

func Validator() string {
	return `
package service

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func errorMessage(e validator.FieldError) map[string]interface{} {
	if e.ActualTag() == "required" {
		return map[string]interface{}{strings.ToLower(string(e.Field())): "not null!"}
	} else if e.ActualTag() == "min" || e.ActualTag() == "max" {
		return map[string]interface{}{strings.ToLower(string(e.Field())): fmt.Sprintf("%v %v char!", strings.ToLower(e.ActualTag()), e.Param())}
	} else if e.ActualTag() == "email" {
		return map[string]interface{}{strings.ToLower(string(e.Field())): "must be email!"}
	} else if e.ActualTag() == "oneof" {
		return map[string]interface{}{strings.ToLower(string(e.Field())): fmt.Sprintf("must be one of: %v", e.Param())}
	} else {
		return map[string]interface{}{strings.ToLower(string(e.Field())): "uups"}
	}

}

func Validator(err error, c *gin.Context) {
	errorMessages := []map[string]interface{}{}
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {

		for _, e := range err.(validator.ValidationErrors) {

			newMessage := errorMessage(e)
			errorMessages = append(errorMessages, newMessage)
		}

	}

	c.JSON(http.StatusBadRequest, errorMessages)
}

func BindingValidator(c *gin.Context, binding interface{}) bool {
	if err := c.BindJSON(binding); err != nil {
		Validator(err, c)
		return true
	}
	return false
}

func BindingValidatorQuery(c *gin.Context, binding interface{}) bool {
	if err := c.BindQuery(binding); err != nil {
		Validator(err, c)
		return true
	}
	return false
}

func MessageResponse(label string, message string) gin.H {
	return gin.H{label: message}
}

	`
}
