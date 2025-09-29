package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/fox998/Chirpy/internal/chirpy_mux"
	"github.com/fox998/Chirpy/internal/config"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config := config.NewApiConfig()

	mux := chirpy_mux.CreateServerMux(config)
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	fmt.Println("Starting to Listen and Serve...")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
}
