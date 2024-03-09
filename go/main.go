package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	srv := New()
	srv.init()

	srv.Mux.Handle("GET /uuid", Recovery(srv.ListUUID))
	srv.Mux.Handle("POST /uuid", Recovery(srv.AddUUID))
	srv.Mux.Handle("GET /ready", Recovery(srv.Readiness))

	addr := fmt.Sprintf("%s:%d", srv.Api.Host, srv.Api.Port)
	log.Printf("running api at %v\n", addr)

	server := &http.Server{Addr: addr, Handler: srv.Mux}
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGSTOP, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("error while shutting down: %v\n", err)
	}

	srv.DB.Close()
}
