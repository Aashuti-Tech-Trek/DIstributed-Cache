package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"distributed-cache/server"
)

func main() {
	port := flag.String("port", "8000", "server port")
	capacity := flag.Int("capacity", 100, "cache capacity")
	flag.Parse()

	self := fmt.Sprintf("127.0.0.1:%s", *port)
	srv := server.NewServer(*capacity, self, []string{}) // no replicas for now

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Distributed Cache is running. Use /get and /set.\n"))
	})
	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)
	http.HandleFunc("/save", srv.SaveHandler)
	http.HandleFunc("/load", srv.LoadHandler)

	fmt.Printf("Server running on :%s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
