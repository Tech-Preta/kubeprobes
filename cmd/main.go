package main

import (
	"log"
	"net/http"

	"github.com/your-package-path/pkg" // Add the import statement for the "pkg" package
)

func main() {
	http.HandleFunc("/checkProbes", pkg.Handler)

	log.Fatal(http.ListenAndServe(":8085", nil))
}
