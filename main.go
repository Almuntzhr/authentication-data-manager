package main

import (
	"context"
	"crypto/tls"
	"log"
	"ms-authentication-data-manager/app"
	"ms-authentication-data-manager/config"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Println("Starting application..")
	config.InitConfig()
	// postgre, err := config.InitPostgreSsh() // for ssh postgree
	postgre, err := config.InitPostgre() // for local postgree
	if err != nil {
		log.Println("Error connect to Postgres: ", err)
	}
	defer postgre.Close()

	router := app.InitRouter(postgre)

	server := &http.Server{
		Addr:         ":" + config.CONFIG["MS_PORT"],
		Handler:      router,
		TLSConfig:    &tls.Config{},
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
		ConnState: func(net.Conn, http.ConnState) {
		},
		ErrorLog: &log.Logger{},
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server.")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting.")
}
