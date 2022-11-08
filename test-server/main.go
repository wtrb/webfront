package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg := fmt.Sprintf("your lucky number is %v", rand.Int())
		fmt.Fprint(w, msg)
	})

	http.ListenAndServe(":8080", nil)
}
