package main

import (
	"log"

	"github.com/Komilov31/sales-tracker/cmd/app"
	_ "github.com/Komilov31/sales-tracker/docs"
)

// @title Sales Tracker API
// @version 1.0
// @description API for tracking sales, income, and expenses with analytics and CSV exports
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	if err := app.Run(); err != nil {
		log.Fatal("could not start server: ", err)
	}
}
