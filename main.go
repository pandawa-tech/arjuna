package main

import (
    "fmt"
    "log"
    "net/http"
	"os"

	"arjuna/middleware"
)

func main() {
	fmt.Fprintf(os.Stdout, "Web Server started. Listening on port 8080\n")
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}