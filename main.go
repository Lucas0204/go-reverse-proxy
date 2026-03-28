package main

import (
	"errors"
	"go-reverse-proxy/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	requestHandler := handler.NewRequestHandler()

	go func() {
		if err := http.ListenAndServe(":3333", requestHandler); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()

	log.Println("Server started on port 3333")

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	log.Println("Server stopped")
}
