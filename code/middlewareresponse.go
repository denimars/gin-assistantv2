package code

func MiddlewareResponse() string {
	return `
package middleware

type message struct {
	Message string ` + "`json:\"message\"`" + `
}
`
}
