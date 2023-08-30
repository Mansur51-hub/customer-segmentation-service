package main

import (
	"github.com/Mansur51-hub/customer-segmentation-service/app"
)

const configPath = "./config"

// @title customer segmentation service
// @version 1.0
// @description swagger
// @termsOfService http://swagger.io/terms/
// @contact.name Mansur Mamedov
// @contact.email mansyr001mamedov@mail.ru
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath
func main() {
	app.Run(configPath)
}
