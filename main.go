package main

import(
	"github.com/villegasl/urlshortener/models"
	"github.com/villegasl/urlshortener/web"
)

func main() {
	// boltdb service
	models.Start()

	// api
	web.Start()
}