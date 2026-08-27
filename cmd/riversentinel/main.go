package main

import (
	"log"
	"net/http"
	"os"

	"github.com/wyw14/cry-153/internal/api"
	"github.com/wyw14/cry-153/internal/service"
)

func main() {
	addr := os.Getenv("RIVERSENTINEL_ADDR")
	if addr == "" {
		addr = ":21253"
	}
	runtime := service.NewRuntime()
	server := &http.Server{Addr: addr, Handler: api.NewRouter(runtime)}
	log.Printf("riversentinel listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
