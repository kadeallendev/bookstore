package main

import (
	"fmt"

	"bookstore/server"
)

func main() {
	server := server.NewServer()
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
