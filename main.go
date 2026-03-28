package main

import (
	"errors"
	"go-reverse-proxy/core"
	"go-reverse-proxy/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	rootDir, _ := os.Getwd()
	configFilePath := rootDir + "/etc/xnign/xnign.conf"
	if !fileExists(configFilePath) {
		panic("Config file not found")
	}

	proxyMiddleware := core.NewProxyMiddleware(configFilePath)
	requestHandler := handler.NewRequestHandler(proxyMiddleware)

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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
