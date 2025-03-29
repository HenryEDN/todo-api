package main

import(
	"log"
)

func main(){
	log.Println("starting server...")
	server := NewAPIServer(":9999")
	server.Run()
}