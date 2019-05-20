package main

import(
	"github.com/villegasl/urlshortener/models"
	"github.com/villegasl/urlshortener/web"
)

func main() {
	// boltdb service
	DB_Handler := models.Start()

	// api
	web.Start(DB_Handler)
}