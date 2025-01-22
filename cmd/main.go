package main

import "github.com/njayp/theseus/pkg/server"

func main() {
	server, err := server.NewServer()
	if err != nil {
		panic(err)
	}

	if err := server.Start(8080); err != nil {
		panic(err)
	}
}
