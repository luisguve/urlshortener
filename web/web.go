package web

import (
	"log"
	"fmt"
	"net/http"
	"time"

	"github.com/villegasl/urlshortener/web/api"
	"github.com/villegasl/urlshortener/web/www"
	"github.com/villegasl/urlshortener/models"

	"github.com/gorilla/mux"
)

func Start() {
	log.Println("web starts")
	router := mux.NewRouter()

	router.HandleFunc("/", www.Index).Methods(http.MethodGet)

	router.Handle("/api/shorturl/{url:[a-zA-Z0-9]{1,11}}", 
		api.RedirectByShortURL()).Methods(http.MethodGet)
	router.HandleFunc("/api/shorturl/new", 
		api.NewShortURL).Methods(http.MethodPost)

	//serve the static files (CSS)
	cssHandler := http.StripPrefix("/static", http.FileServer(http.Dir("./static")))
	router.PathPrefix("/static/").Handler(cssHandler)

	srv := http.Server {
		Addr: 			":8080",
		Handler: 		router,
		ReadTimeout: 	10 * time.Second,
		WriteTimeout: 	10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("visit localhost:8080")
	log.Fatal(srv.ListenAndServe())
}