package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/koyokr/sjb-host/api"
	"github.com/koyokr/sjb-host/db"
	"github.com/rs/cors"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := db.InitDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseDB()

	ch := make(chan string)

	go func() {
		for {
			updateDomainAll()
		}
	}()
	go func() {
		for {
			insertDomainFromChan(ch)
		}
	}()

	router := httprouter.New()
	router.GET("/", api.GetDomainAll)
	router.GET("/:name", api.PutDomainFunc(ch))

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
