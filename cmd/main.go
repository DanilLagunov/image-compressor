package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DanilLagunov/image-compressor/pkg/api"
	localstorage "github.com/DanilLagunov/image-compressor/pkg/storage/local_storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage := localstorage.New("././images/")
	server := http.Server{
		Addr:              ":8080",
		Handler:           api.New(ctx, storage),
		ReadHeaderTimeout: 60 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
	}

	fmt.Println("Listening")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
