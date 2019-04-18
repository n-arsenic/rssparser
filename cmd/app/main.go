package main

import (
	"fmt"
	"log"
	"net/http"
	"rssparser/internal/pkg/handlers"
)

func main() {
	fmt.Println("Server listen...")
	handlers.Compose()
	log.Fatal(http.ListenAndServe(":8282", nil))
}
