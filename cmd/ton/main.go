package main

import (
	"log"
	"net/http"

	"github.com/ton/framework/internal/orchestrator"
)

func main() {
	o := orchestrator.NewOrchestrator()
	http.HandleFunc("/tools", o.ServeHTTP)
	log.Println("TON server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}