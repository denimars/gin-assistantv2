package code

func Main(project string) string {
	return `
package main

import "` + project + `/app"

func main() {
	app.Run()
}
	`
}
