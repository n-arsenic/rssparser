package main

import (
	"fmt"
	"log"
	"net/http"
	"rssnews/api/handlers"
)

func main() {
	fmt.Println("Server listen...")
	handlers.Compose()
	log.Fatal(http.ListenAndServe(":8282", nil))
}
