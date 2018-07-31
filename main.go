package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/koyokr/sjb-host/api"
	"github.com/rs/cors"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := api.Init(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer api.Close()

	router := httprouter.New()
	router.GET("/", api.GetDomains)

	router.NotFound = http.HandlerFunc(api.NotFound)
	router.MethodNotAllowed = http.HandlerFunc(api.MethodNotAllowed)

	handler := cors.Default().Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	log.Println("Strating HTTP server on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
