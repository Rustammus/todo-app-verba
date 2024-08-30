package main

import "ToDoVerba/internal/app"

// @title           ToDo service
// @version         1.0
// @description     This is my server.

// @license.name  Apache helicopter
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8082
// @BasePath  /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	app.Run()
}
