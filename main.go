package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/products", httpHandler)
	fmt.Println("Server has started ✅")
	http.ListenAndServe(":8080", nil)
}

func httpHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	response, err := getProducts()

	if err != nil {
		fmt.Fprintf(w, err.Error()+"\r\n")
	} else {
		enc := json.NewEncoder(w)
		enc.SetIndent("", " ")

		if err := enc.Encode(response); err != nil {
			fmt.Println(err.Error())
		}
	}

}
