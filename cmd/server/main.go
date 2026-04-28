package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "DEBUG: El servidor esta vivo. ID: 9999")
	})

	// ESTE MENSAJE DEBE APARECER EN RENDER
	fmt.Printf("🚀 REVISIÓN TÉCNICA 9999 - Puerto %s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
